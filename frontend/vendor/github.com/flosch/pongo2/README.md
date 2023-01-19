# [pongo](https://en.wikipedia.org/wiki/Pongo_%28genus%29)2

[![GoDoc](https://godoc.org/github.com/flosch/pongo2?status.png)](https://godoc.org/github.com/flosch/pongo2)
[![Build Status](https://travis-ci.org/flosch/pongo2.svg?branch=master)](https://travis-ci.org/flosch/pongo2)
[![Coverage Status](https://coveralls.io/repos/flosch/pongo2/badge.png?branch=master)](https://coveralls.io/r/flosch/pongo2?branch=master)
[![gratipay](http://img.shields.io/badge/gratipay-support%20pongo-brightgreen.svg)](https://gratipay.com/flosch/)
[![Bountysource](https://www.bountysource.com/badge/tracker?tracker_id=3654947)](https://www.bountysource.com/trackers/3654947-pongo2?utm_source=3654947&utm_medium=shield&utm_campaign=TRACKER_BADGE)

pongo2 is the successor of [pongo](https://github.com/flosch/pongo), a Django-syntax like templating-language.

Install/update using `go get` (no dependencies required by pongo2):
```
go get -u github.com/flosch/pongo2
```

Please use the [issue tracker](https://github.com/flosch/pongo2/issues) if you're encountering any problems with pongo2 or if you need help with implementing tags or filters ([create a ticket!](https://github.com/flosch/pongo2/issues/new)). If possible, please use [playground](https://www.florian-schlachter.de/pongo2/) to create a short test case on what's wrong and include the link to the snippet in your issue.

**New**: [Try pongo2 out in the pongo2 playground.](https://www.florian-schlachter.de/pongo2/)

## First impression of a template

```HTML+Django
<html><head><title>Our admins and users</title></head>
{# This is a short example to give you a quick overview of pongo2's syntax. #}

{% macro user_details(user, is_admin=false) %}
	<div class="user_item">
		<!-- Let's indicate a user's good karma -->
		<h2 {% if (user.karma >= 40) || (user.karma > calc_avg_karma(userlist)+5) %}
			class="karma-good"{% endif %}>
			
			<!-- This will ca