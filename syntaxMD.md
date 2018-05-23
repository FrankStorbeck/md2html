Mastering Markdown
==================

## Headers
# This is an &lt;h1&gt; tag
## This is an &lt;h2&gt; tag
###### This is an &lt;h6&gt; tag

##Styling text
normal **bold** normal

normal __bold__ normal

normal *italic* normal

normal ~~deleted~~ normal

**This text is _extremely_ important ~~deleted~~ text**

##Ignoring Markdown formatting
normal \~~normal also\~~ normal

normal \*normal also\* normal

normal \**italic(!)\** normal

normal \_normal also\_ normal

normal \__italic(!)\__ normal

##Quoting text
> quote 1
>quote 2

##Quoting inline code
Here is some `inline` code.

##Quoting code
```
multi
line
code
```
or
    multi
    line
    code

##Links
Link to [something](https://guides.github.com/features/mastering-markdown/),

or to something [local](doc/README.md).

Can also link to images: ![an image](img/img.png).

##Lists
###Unordered
* Unordered list item 1
* Unordered list item 2
* Unordered list item 3
  - Unordered list item 3.1
  - Unordered list item 3.2
    * Unordered list item 3.2.1
    * Unordered list item 3.2.2
  - Unorderd list item 3.3
    * Unorderd list item 3.3.1
* Unorderd list item 4

###Ordered
1. Ordered list item 1
2. Ordered list item 2
3. Ordered list item 3
   1. Ordered list item 3.1
   2. Ordered list item 3.2

###Mixed ordered and unordered
* Unordered list item 1
* Unordered list item 2
* Unordered list item 3
  3. Ordered list item 3.1
  3. Ordered list item 3.2
     * Unordered list item 3.2.1
     * Unordered list item 3.2.2
  3. Orderd list item 3.3
     * Unorderd list item 3.3.1
* Unorderd list item 4

##Organizing information with tables
| AAA | NNN  | MMMMM | WWWWWWW |
| --- | :--- | ----: | :-----: |
| aaa | nnn  | mmmmm | \|wwwwww |
| zzz |
| 000 | 111  | 22222 | 3333333 | 44444 |
