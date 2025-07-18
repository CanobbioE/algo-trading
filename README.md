# algo-trading

Using algorithms to make informed trades.

⚠️
> This project was created as a fun learning experience and personal challenge.  
> It's not guarantee to be mantained.  
> It does not provide professional insight on stock trading.

⚠️

---

This project it delivered as a go CLI tool, check the at [Available Commands](./docs/available-commands.md).

All commands can be configured using a JSON file.
Check the documentation about the tool configuration [here](./docs/configuration.md)  

There is the option to use a LLM model to analyse the output of some commands and provide insight as a virtual "AI" assistant.
To configure the "AI" assistant, a JSON config file is required with the following structure:
```json
{
  "api_key": "xxxx"
}
```


## Run

### Building from source
```text
go build .
```

```text
./algo-trading [command] [flags]
```

### Run without building a binary

```text
go run . [command] [flags]
```
