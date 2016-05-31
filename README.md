[![wercker status](https://app.wercker.com/status/701d88d0a7128c919abdbf9935a1d632/m "wercker status")](https://app.wercker.com/project/bykey/701d88d0a7128c919abdbf9935a1d632)

# World of FluxCraft - Reality
This is the reality (game state) service for the World of FluxCraft sample. It is responsible solely for allowing for the update (PUT) and query (GET) of completed/cached game state. The event processor is responsible for computing the game state. This service is a cache that is used by the UI to load the state of existing games when the Flux application initializes.
