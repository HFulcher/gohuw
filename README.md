# gohuw
A Golang static site generator.

## Introduction
This has not been written necessarily for public use. I wrote it to learn more about Go and also challenge myself to write an SSG. Feel free to use it but be aware:

1. I am unlikely to provide support in the case there are issues.
2. This is a very opinionated SSG, it is tailored to my needs/wants.
3. There is no guarantee this will work in the long run.
4. It is probably (more like definitely) not using Go in the best way.

## Usage
Drop the binary into the folder containing your website and either run:

- `./gohuw` - Compile your folder into a static website in `public/`
- `./gohuw dev` - Run a localhost server that serves your website from `public/`. This will rebuild when any file in the folder changes.
    - `-p` (Defaults to `8100`) - Specify the port to serve from
    - `-s` (Default to `public/`) - Directory to serve files from
    - `-w` (Default to `.`) - Directory to watch for changes from

It should be structured in the following way:

```
root
|- content/ 
|  |- blog/
|  |  |- index.md 
|  |  |- my-blog-post.md
|  |  |- my-other-blog-post.md
|  |- essays/
|- templates/ 
|  |- baseof.html
|  |- single.html
|  |- list.html
|  |- other_template.html
|- assets/ 
|- public/ 
|- index.md 
|- any_other_page.md 
|- config.json 
|- gohuw
```

### `content/` 
This is where all collections of .md files should go. Each collection type should be in its own folder with an `index.md` file that acts as the "list" page.

#### `index.md`
The `index.md` template will have access to a `.Pages` variable that will provide a list of the other .md files in reverse chronological order.

Frontmatter for an `index.md` should include `layout: <some layout name>` where the layout name corresponds to the name of a template in `templates` (without the extension). Otherwise `list.html` will be used by default.

#### `.md` 
For a markdown collection item the only frontmatter needed is a `date` field to be used in sorting.


### `templates/`
This is where you store all your template files that use [html/template](https://pkg.go.dev/html/template). The only required files in this folder are:

- `baseof.html` - This is your "layout" file that is applied to all other templates.
- `single.html` - This is the default template used if a `layout` frontmatter isn't provided.
- `list.html` - This is the default template used only for `index.md` files in `content/` if the `layout` frontmatter isn't provided.

To use your own templates you can specify the name of the file (minus extension) in the `layout` frontmatter of a file.

You have access to a number of useful items in your template:

- `.Site` - This is all the information found in your `config.json`. The only default items in `.Site` are:
    - `.Site.Url` - In dev mode this will default to `localhost:<port specified>`. If not in dev mode this needs to be set in `config.json` as the base Url.
    - `.Site.IsProduction` - A boolean that is true when `./gohuw` is run and false when `./gohuw dev` is run. Useful for filtering out production only pieces of template (e.g. analytics trackers).
- `.Page` - This contains information and content about the current page.
    - `.Page.Title` - The title of the markdown file specified in the `title` frontmatter field. If this field is not present `.Page.Title` defaults to an empty string.
    - `.Page.Slug` - The Url path of the page. This does not include the base part of the Url.
    - `.Page.Path` - The filepath of the markdown file.
    - `.Page.Destination` - The filepath of the converted HTML file.
    - `.Page.Content` - The contents of the markdown file converted to HTML. This is primarily what you will use to populate your HTML from the templates.
    - `.Page.Metadata` - All fields from the frontmatter. These are accessible as follows `.Page.Metadata['field_name'].
- `.Pages` (`content/x/index.md` files only) - This provides a list of Markdown files from the relevant content folder in reverse chronological order. This is useful for producing list pages.


### `assets/`
This is where your static resources should be stored. This folder will be copied as is to a `static/` folder in `public/`.


### `public/`
This is where your compiled files will be stored.


### `config.json`
This is where you can configure any site wide variables that you might wants (accessible by `.Site` in the templates). There are no required variables in this file but it must be present for gohuw to run.

### Top level `.md` files
Files stored at the top level of the folder will be treated as single pages to be converted to html. This is for things such as your homepage `index.md` or contact pages.


## Deployment
Depending on how you deploy your website you have a few options. For SSG hosts like Netlify you can provide your repository as a folder and point Netlify to serve from the `public/` folder.

If you want to generate your website as part of the build process you can either bundle the `gohuw` binary or `curl` it at build time from the release links.