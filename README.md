# reflex

A high level implementation of reflecting and setting data types, ideal for manipulating custom structs. This package was developed in support of another package that I'm currently developing. If you're looking for a more general purpose and friendly way of reflecting types and values then I suggest you check out [oleiade/reflections](https://github.com/oleiade/reflections).

## Installation
`go get github.com/akillmer/reflex`

## Setting a field
```go
var s = struct {
        Name struct {
            First string
            Last  string
        }
    Location string
}{}

var m = reflex.NewModel(&s)

m.Set("Name.First", "Andrew")
m.Set("Name.Last", "Killmer")
m.Set("Location", "Aiea, HI")
```

## Mapping fields
```go
var data = make(map[string]interface{})
data["Name.First"] = "Andrew"
data["Name.Last"] = "Killmer"
// etc

m.Map(data)
```

The key difference between `Map` and `Set` is when it comes down to slices. You may want to populate multiple fields of a struct before adding it to the model's slice. Using `Set` would result in a new element being added for every time it's called.

## All setter, no getter
Since `reflex` needs an `interface{}` to work then obviously it doesn't need to provide any getter functions as you have access to the original struct. 

## Caveats
Naturally, making sure the value you're assigning to the model is of the same type. If not then your program will panic. Also, I have not implemented support for `map` types yet.