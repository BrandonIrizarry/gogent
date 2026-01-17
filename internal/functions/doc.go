/*
Package functions serves as the bridge between the model's intent and
the guardrailed set of system capabilities.

It is responsible for:

 1. Defining the schemas for available functions.

 2. Validating incoming function call requests from the LLM.

 3. Executing the actual business logic associated with each function call.

 4. Formatting the results back into a structure consumable by the LLM.
*/
package functions
