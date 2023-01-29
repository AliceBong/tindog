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
			
			<!-- This will call user.String() automatically if available: -->
			{{ user }}
		</h2>

		<!-- Will print a human-readable time duration like "3 weeks ago" -->
		<p>This user registered {{ user.register_date|naturaltime }}.</p>
		
		<!-- Let's allow the users to write down their biography using markdown;
		     we will only show the first 15 words as a preview -->
		<p>The user's biography:</p>
		<p>{{ user.biography|markdown|truncatewords_html:15 }}
			<a href="/user/{{ user.id }}/">read more</a></p>
		
		{% if is_admin %}<p>This user is an admin!</p>{% endif %}
	</div>
{% endmacro %}

<body>
	<!-- Make use of the macro defined above to avoid repetitive HTML code
	     since we want to use the same code for admins AND members -->
	
	<h1>Our admins</h1>
	{% for admin in adminlist %}
		{{ user_details(admin, true) }}
	{% endfor %}
	
	<h1>Our members</h1>
	{% for user in userlist %}
		{{ user_details(user) }}
	{% endfor %}
</body>
</html>
```

## Development status

**Latest stable release**: v3.0 (`go get -u gopkg.in/flosch/pongo2.v3` / [`v3`](https://github.com/flosch/pongo2/tree/v3)-branch) [[read the announcement](https://www.florian-schlachter.de/post/pongo2-v3/)]

**Current development**: v4 (`master`-branch)

*Note*: With the release of pongo v4 the branch v2 will be deprecated.

**Deprecated versions** (not supported anymore): v1

| Topic                                | Status                                                                                 |
| ------------------------------------ | -----------------------------------------------------------------------------------