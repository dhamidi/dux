# Description

Dux is a code generator for you project.  It allows you to easily generate new components in your application according to a common blueprint.  By using Dux you interact with your project not on a file-based level, but on a semantic level.  Which of the following do you prefer?

*Without using Dux*

```sh
$ touch lib/my_app/controllers/widgets_controller.rb
$ $EDITOR lib/my_app/controllers/widgets_controller.rb
# write controller boilerplate
$ touch test/my_app/controllers/widgets_controller_test.rb
# write controller test boilerplate
```

*Using Dux*

```sh
$ dux new controller Widgets
create lib/my_app/controllers/widgets_controller.rb
create test/my_app/controllers/widgets_controller.rb
$ dux edit controller Widgets
$EDITOR lib/my_app/controllers/widgets_controller.rb
```

# Why Dux

Dux helps you keep your project structured and clean.  It works in projects of any programming language, helping you build your own conventions.

Ruby on Rails popularized the idea of "convention over configuration" and comes with [tools to build new conventions](http://guides.rubyonrails.org/generators.html) in your own Rails project.  Other frameworks have followed suit and provide similar functionality.  However, every project is different and so are the conventions within this project.

If it was simple and easy to encode new conventions using a tool, new team members can get up to speed more quickly by having a visible representation of those conventions.

# How it works

At its core Dux is a file generator based on text templates.  A set of templates is called a *blueprint*.  Dux creates files for a blueprint in three steps:

1. *Gather data*: Dux gathers data about your project from various sources and makes them available in the blueprint's templates.  External programs can be used in order to get deep insights into your project, such as which modules, functions and classes are defined already.
2. *Render blueprint templates*: Templates files are rendered into a staging area, together with the data that has been gathered in step 1.  Any errors that occurred during template rendering are logged and the templates and data can be inspected in the staging are in case there were any problems.
3. *Copy rendered files into your project*: The files generated in step 2 are now copied into your project's directory tree. Dux will not overwrite any existing files when copying but warn you in that case.
4. *Edit existing files*: if a blueprint specifies any editing operations on files, these operations are now executed on

Using text templates allows Dux to stay flexible enough to work with any programming language.
