![CI Status](https://github.com/matt-bourke/mdx/actions/workflows/go.yml/badge.svg)

# MDX
MDX is a custom prototype markdown extension format where custom data can be added to regular markdown tags.
It aims to give a little extra control over the structure and function of markdown files by allowing you to add 
a variety of elements. The main focus of MDX is for converting content into HTML format.

It's still pretty (very) rough around the edges. See the files `sample.mdx` and `sample.html` for some examples.

MDX can be used as a command line tool, or as an imported package.


### Quick Links
[Extensions](#Extensions)
[Usage](#Usage)


## Extensions
### Properties
To add more customisability to markdown, MDX features properties. By prefixing elements with name/value properties
wrapped in `{ }`, the subsequent parsed elements will receive these properties when parsed into HTML.

Example:
```md
{ .class=section-heading }
# Welcome
```

### Divs
To add more structure, divs can be parsed into the HTML by wrapping content in `[ ]`. Combining divs with properties
allows for much more control over the styling and structure of the resulting HTML.

Example:
```md
[
  # Welcome
  Hello there
]
```

### Spans
For inline structure, spans can be parsed by wrapping content in `$ $`. Combining spans with properties is also a
powerful way of managing inline styling.

Example:
```md
Hello, $world$!
```

### Buttons
For adding interactivity, buttons can also be added with the syntax `~[x](y)`, where x is the button label, and y is
the name of the click handler function defined in the HTML script tag. Still figuring out a good way to get code in the
script tags.

Example:
```md
~[Click Me](handleClick)
```

### Custom Code Block
This one generates some very specific styling for a particular use case, and is the catalsyst for MDX being created.
To generate the custom code block, wrap the code in `^^ ^^`. Syntax highlighting is not supported but hopefully will be
in the future.

Example:
```md
^^
func main() {
  fmt.Println("Hello, world!")
}
^^
```

### Nav
While not entirely useful since we already have divs and custom properties, Nav elements are supported anyways. You can
generate a Nav element by wrapping the content in `@ @`.

Example:
```md
@
[Home](/home)
[Feed](/feed)
[Account](/account)
@
```

### Comments
Yes, markdown already supports comments, but I prefer commenting a line by prefixing it with `//`, so that's what I've
done here.

Example:
```md
// do you really need an example for this one?
```


## Usage
### Import
After importing MDX into your Go project, you can generate HTML files by calling the `Generate` method, and passing it
an appropriate `config` object. The config option must have `Filename` data, and can optionally have other data such as:
- Title
- OutputFilename
- CssFilename
- Favicon
- FontLink

### Command Line
TODO :P
