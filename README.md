md2html
=======

golang program to translate the mark down code into an html
version.

Usage
-----
```
> md2html -help
Usage of md2html:
  -in string
    	path to input file (default "stdin")
  -out string
    	path to output file (default "stdout")
  -style string
    	style sheet for HTML document
  -title string
    	title for HTML document
  -version
    	Show version number and exit
```

Recognised formatting
---------------------

hd2htm adheres to
https://daringfireball.net/projects/markdown/syntax

but note:

- `styling` and `quoting` a text must be terminated before the end of a line.
- a `link` can be made to the start of a `header` by replacing all spaces in it
   by a hyphen and all capitals by it lowercase equivalent.
- `list` items can also start with a '+' character.
- when a `list` item contains an empty line all entries for that item are put
  into a paragraph.
- task lists are not implemented.


See also:
https://guides.github.com/features/mastering-markdown/
or
https://help.github.com/articles/basic-writing-and-formatting-syntax/
