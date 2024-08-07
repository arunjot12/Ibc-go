<!--
order: 6
-->

# Handling Updates and Misbehaviour

The functions for handling updates to a light client and evidence of misbehaviour are all found in the [`ClientState`](https://github.com/cosmos/ibc-go/blob/v6.0.0/modules/core/exported/client.go#L40) interface, and will be discussed below.

It is important to note that `Misbehaviour` in this particular context is referring to misbehaviour on the chain level intended to fool the light client. This will be defined by each light client. 

## `VerifyClientMessage` 

`VerifyClientMessage` must verify a `ClientMessage`. A `ClientMessage` could be a `Header`, `Misbehaviour`, or batch update. To understand how to implement a `ClientMessage`, please refer to the [documentation](./update.md)

It must handle each type of `ClientMessage` appropriately. Calls to `CheckForMisbehaviour`, `UpdateState`, and `UpdateStateOnMisbehaviour` will assume that the content of the `ClientMessage` has been verified and can be trusted. An error should be returned if the `ClientMessage` fails to verify.

For an example of a `VerifyClientMessage` implementation, please check the [Tendermint light client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint/update.go#L20).

## `CheckForMisbehaviour`

Checks for evidence of a misbehaviour in `Header` or `Misbehaviour` type. It assumes the `ClientMessage` has already been verified.

For an example of a `CheckForMisbehaviour` implementation, please check the [Tendermint light client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint/misbehaviour_handle.go#L18).

> The Tendermint light client [defines `Misbehaviour`](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint/misbehaviour.go) as two different types of situations: a situation where two conflicting `Headers` with the same height have been submitted to update a client's `ConsensusState` within the same trusting period, or that the two conflicting `Headers` have been submitted at different heights but the consensus states are not in the correct monotonic time ordering (BFT time violation). More explicitly, updating to a new height must have a timestamp greater than the previous consensus state, or, if inserting a consensus at a past height, then time must be less than those heights which come after and greater than heights which come before.

## `UpdateStateOnMisbehaviour`

`UpdateStateOnMisbehaviour` should perform appropriate state changes on a client state given that misbehaviour has been detected and verified. This method should only be called when misbehaviour is detected, as it does not perform any misbehaviour checks. Notably, it should freeze the client so that calling the `Status()` function on the associated client state no longer returns `Active`.

For an example of a `UpdateStateOnMisbehaviour` implementation, please check the [Tendermint light client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint/update.go#L197).

## `UpdateState`

`UpdateState` updates and stores as necessary any associated information for an IBC client, such as the `ClientState` and corresponding `ConsensusState`. It should perform a no-op on duplicate updates.

It assumes the `ClientMessage` has already been verified.

For an example of a `UpdateState` implementation, please check the [Tendermint light client](https://github.com/cosmos/ibc-go/blob/main/modules/light-clients/07-tendermint/update.go#L131).

## Putting it all together

The `02-client Keeper` module in ibc-go offers a reference as to how these functions will be used to [update the client](https://github.com/cosmos/ibc-go/blob/main/modules/core/02-client/keeper/client.go#L48).

```golang
if err := clientState.VerifyClientMessage(clientMessage); err != nil {
        return err
    }
    
    foundMisbehaviour := clientState.CheckForMisbehaviour(clientMessage)
    if foundMisbehaviour {
        clientState.UpdateStateOnMisbehaviour(clientMessage)
        // emit misbehaviour event
        return 
    }
    
    clientState.UpdateState(clientMessage) // expects no-op on duplicate header
    // emit update event
    return
}

