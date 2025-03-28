---
Styles:
  - /static/styles/code.css
Scripts:
  - /static/scripts/copy-code.js
---

# Markdown Examples

## Headings

    # h1 Heading
    ## h2 Heading
    ### h3 Heading
    #### h4 Heading
    ##### h5 Heading
    ###### h6 Heading

# h1 Heading
## h2 Heading
### h3 Heading
#### h4 Heading
##### h5 Heading
###### h6 Heading

## Horizontal Rules

    ___

    ---

    ***

___

---

***

## Typographic replacements

...  -- ---

"Smartypants, double quotes" and 'single quotes'

## Emphasis

    **This is bold text**

    __This is bold text__

    *This is italic text*

    _This is italic text_

    ~~Strikethrough~~

**This is bold text**

__This is bold text__

*This is italic text*

_This is italic text_

~~Strikethrough~~

## Blockquotes

    > Blockquotes can also be nested...
    >> ...by using additional greater-than signs right next to each other...
    > > > ...or with spaces between arrows.

> Blockquotes can also be nested...
>> ...by using additional greater-than signs right next to each other...
> > > ...or with spaces between arrows.

## Lists

### Unordered

    + Create a list by starting a line with `+`, `-`, or `*`
    + Sub-lists are made by indenting 2 spaces:
      - Marker character change forces new list start:
        * Ac tristique libero volutpat at
        + Facilisis in pretium nisl aliquet
        - Nulla volutpat aliquam velit
    + Very easy!

+ Create a list by starting a line with `+`, `-`, or `*`
+ Sub-lists are made by indenting 2 spaces:
  - Marker character change forces new list start:
    * Ac tristique libero volutpat at
    + Facilisis in pretium nisl aliquet
    - Nulla volutpat aliquam velit
+ Very easy!

### Ordered

    1. Lorem ipsum dolor sit amet
    2. Consectetur adipiscing elit
    3. Integer molestie lorem at massa

1. Lorem ipsum dolor sit amet
2. Consectetur adipiscing elit
3. Integer molestie lorem at massa

#### Start numbering with offset:

    42. foo
    43. bar

42. foo
42. bar

## Code

    Inline code: `Inline code`

Inline code: `Inline code`

### Indented code

    // Some comments
    line 1 of code
    line 2 of code
    line 3 of code

### Block code "fences"

    ```
    Sample text here...
    ```

```
Sample text here...
```

#### Block code "fences" with syntax highlighting

  	``` go
  	var test string
  	test = "Test"
  	fmt.Println(test)

  	anotherTest := "Test2"
  	fmt.Println(anotherTest)
  	functionCallWithLotsOfParams(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, test)
  	```

``` go
var test string
test = "Test"
fmt.Println(test)

anotherTest := "Test2"
fmt.Println(anotherTest)
functionCallWithLotsOfParams(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, nil, test)
```

## Tables

    | Option | Description |
    | ------ | ----------- |
    | data   | path to data files to supply the data that will be passed into templates. |
    | engine | engine to be used for processing templates. Handlebars is the default. |
    | ext    | extension to be used for dest files. |

| Option | Description |
| ------ | ----------- |
| data   | path to data files to supply the data that will be passed into templates. |
| engine | engine to be used for processing templates. Handlebars is the default. |
| ext    | extension to be used for dest files. |

### Right aligned columns

    | Option | Description |
    | ------:| -----------:|
    | data   | path to data files to supply the data that will be passed into templates. |
    | engine | engine to be used for processing templates. Handlebars is the default. |
    | ext    | extension to be used for dest files. |

| Option | Description |
| ------:| -----------:|
| data   | path to data files to supply the data that will be passed into templates. |
| engine | engine to be used for processing templates. Handlebars is the default. |
| ext    | extension to be used for dest files. |

## Links

    [link text](http://dev.nodeca.com)

[link text](http://dev.nodeca.com)

    [link with title](http://nodeca.github.io/pica/demo/ "title text!")

[link with title](http://nodeca.github.io/pica/demo/ "title text!")

    Autoconverted link https://github.com/nodeca/pica (enable linkify to see)

Autoconverted link https://github.com/nodeca/pica (enable linkify to see)

## Images

    ![Minion](https://octodex.github.com/images/minion.png)

![Minion](https://octodex.github.com/images/minion.png)

    ![Stormtroopocat](https://octodex.github.com/images/stormtroopocat.jpg "The Stormtroopocat")

![Stormtroopocat](https://octodex.github.com/images/stormtroopocat.jpg "The Stormtroopocat")

Like links, Images also have a footnote style syntax
With a reference later in the document defining the URL location:

    ![Alt text][id]

    [id]: https://octodex.github.com/images/dojocat.jpg  "The Dojocat"

![Alt text][id]

[id]: https://octodex.github.com/images/dojocat.jpg  "The Dojocat"
