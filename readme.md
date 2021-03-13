# queryoptions

## Query string parameters converted to structure queries

### Filters

At a high level, filters can be grouped as mandatory or optional. Supported filter types include:

* beginsWith
* contains
* endsWith
* exact
* greaterThan
* greaterThanEqual
* lessThan
* lessThanEqual
* notEqual

The filters are provided via bracketed notation.

`?filter[fieldName]=value`

`?filter[mandatory][exact][fieldName]=value`

```json
{
  "mandatory": {
    "exact": {
      "fieldName": "value"
    }
  }
}
```

`?filter[mandatory][beginsWith][fieldName2]=value&filter[mandatory][exact][fieldName1]=test`

```json
{
  "mandatory": {
    "beginsWith": {
      "fieldName2": "value"
    },
    "exact": {
      "fieldName1": "test"
    }
  }
}
```

`?filter[mandatory][beginsWith][fieldName2]=value&filter[mandatory][exact][fieldName1]=test&filter[optional][exact][someID]=1234`

```json
{
  "mandatory": {
    "beginsWith": {
      "fieldName2": "value"
    },
    "exact": {
      "fieldName1": "test"
    }
  },
  "optional": {
    "exact": {
      "someID": 1234
    }
  }
}
```