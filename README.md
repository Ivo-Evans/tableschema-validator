# Tableschema validator

> Note: This is a learning project and not ready for production use

This is a Go package that constructs a Datapackage v2 schema as defined [here](https://datapackage.org/standard/table-schema/) (UI) and [here](https://datapackage.org/profiles/2.0/tableschema.json) (JSON).

Limitations include:
- Only supporting a limited subset of the schema (selected properties of String, Number, Boolean and List fields)
- Not integrating with other parts of the Datapackage standard

The aim is for it to construct a _valid_ tableschema and to validate data against it, but there are still many parts of the spec that will not be implemented. 