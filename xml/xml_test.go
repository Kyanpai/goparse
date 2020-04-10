package xml

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var simpleXml = []byte(`
<root>
    <p>content 0</p>
    <foo>
        <p>foo content 0</p>
        <bar>
           <p>john</p>
           <p>doe</p>
        </bar>
    </foo>
    <baz>
        <p class="baz">baz content 0</p>
        <p>baz content 1</p>
    </baz>
    <foo>
        <p>foo content 1</p>
    </foo>
</root>
`)

var complexeXml = []byte(`
<?xml version="1.0" encoding="UTF-8"?>
<feed xmlns="http://www.w3.org/2005/Atom" xmlns:media="http://search.yahoo.com/mrss/" xml:lang="en-US">
  <id>tag:github.com,2008:https://github.com/NeverSinkDev/NeverSink-Filter/releases</id>
  <link type="text/html" rel="alternate" href="https://github.com/NeverSinkDev/NeverSink-Filter/releases"/>
  <link type="application/atom+xml" rel="self" href="https://github.com/NeverSinkDev/NeverSink-Filter/releases.atom"/>
  <title>Release notes from NeverSink-Filter</title>
  <updated>2020-03-27T22:23:45+01:00</updated>
  <entry>
    <id>tag:github.com,2008:Repository/59748131/7.8.2</id>
    <updated>2020-03-27T22:25:57+01:00</updated>
    <link rel="alternate" type="text/html" href="https://github.com/NeverSinkDev/NeverSink-Filter/releases/tag/7.8.2"/>
    <title>NeverSink&#39;s itemfilter - version 7.8.2 - Major improvement to higher strictnesses and Economy Update #2</title>
    <content type="html">&lt;ul&gt;
&lt;li&gt;
&lt;p&gt;[TIERING] Most tierlists have been adjusted to better match the current state of the game. Affected tierlists are: uniques, divination cards, unique maps, fossils, incubators, prophecies, scarabs, high level crafting bases, oils, fragments, shaper, elder, crusader, hunter, redeemer, warlord, vials, currencies, deliriumorbs&lt;/p&gt;
&lt;li&gt;
&lt;p&gt;Low Quality gems are now hidden on very strict, instead of uber-strict.&lt;/p&gt;
&lt;/li&gt;
&lt;/ul&gt;</content>
    <author>
      <name>NeverSinkDev</name>
    </author>
    <media:thumbnail height="30" width="30" url="https://avatars1.githubusercontent.com/u/2942999?s=60&amp;v=4"/>
  </entry>
  <entry>
    <id>tag:github.com,2008:Repository/59748131/7.8.1</id>
    <updated>2020-03-20T16:51:25+01:00</updated>
    <link rel="alternate" type="text/html" href="https://github.com/NeverSinkDev/NeverSink-Filter/releases/tag/7.8.1"/>
    <title>NeverSink&#39;s itemfilter - version 7.8.1 - Delirium Economy Update #1</title>
    <content type="html">&lt;ul&gt;
&lt;li&gt;
&lt;p&gt;[TIERING] Most tierlists have been adjusted to better match the current state of the game. Affected tierlists are: uniques, divination cards, unique maps, fossils, incubators, prophecies, scarabs, high level crafting bases, oils, fragments, shaper, elder, crusader, hunter, redeemer, warlord, vials, currencies&lt;/p&gt;
&lt;/li&gt;
&lt;/ul&gt;</content>
    <author>
      <name>NeverSinkDev</name>
    </author>
    <media:thumbnail height="30" width="30" url="https://avatars1.githubusercontent.com/u/2942999?s=60&amp;v=4"/>
  </entry>
</feed>

`)

func TestParse(t *testing.T) {

	for _, d := range []struct {
		name string
		xml  []byte
		path string
		ret  []string
	}{
		{
			name: "depth 0",
			xml:  simpleXml,
			path: "root",
			ret:  []string{"content 0"},
		},
		{
			name: "depth 1",
			xml:  simpleXml,
			path: "root/baz",
			ret:  []string{"baz content 0", "baz content 1"},
		},
		{
			name: "depth 1 with index",
			xml:  simpleXml,
			path: "root/foo[1]",
			ret:  []string{"foo content 1"},
		},
		{
			name: "depth 2",
			xml:  simpleXml,
			path: "root/foo[0]/bar",
			ret:  []string{"john", "doe"},
		},
		{
			name: "bad path format",
			xml:  simpleXml,
			path: "root/foo[0]/",
			ret:  []string{},
		},
		{
			name: "complexe xml",
			xml:  complexeXml,
			path: "feed/entry[0]/updated",
			ret:  []string{"2020-03-27T22:25:57+01:00"},
		},
		{
			name: "complexe xml with attributes",
			xml:  complexeXml,
			path: "feed/entry[1]/link@href",
			ret:  []string{"https://github.com/NeverSinkDev/NeverSink-Filter/releases/tag/7.8.1"},
		},
		{
			name: "complexe xml with depth and attributes",
			xml:  complexeXml,
			path: "feed/link[1]@href",
			ret:  []string{"https://github.com/NeverSinkDev/NeverSink-Filter/releases.atom"},
		},
	} {
		t.Run(d.name, func(t *testing.T) {
			ret := Parse(d.path, d.xml)
			assert.Equal(t, d.ret, ret)
		})
	}
}

func TestParseRecursive(t *testing.T) {

	for _, d := range []struct {
		name string
		xml  []byte
		path string
		ret  []string
	}{
		{
			name: "depth 0",
			xml:  simpleXml,
			path: "root",
			ret:  []string{"content 0", "foo content 0", "john", "doe", "baz content 0", "baz content 1", "foo content 1"},
		},
		{
			name: "depth 1",
			xml:  simpleXml,
			path: "root/baz",
			ret:  []string{"baz content 0", "baz content 1"},
		},
		{
			name: "depth 1 with index",
			xml:  simpleXml,
			path: "root/foo[1]",
			ret:  []string{"foo content 1"},
		},
		{
			name: "depth 2",
			xml:  simpleXml,
			path: "root/foo[0]/bar",
			ret:  []string{"john", "doe"},
		},
		{
			name: "bad path format",
			xml:  simpleXml,
			path: "root/foo[0]/",
			ret:  []string{},
		},
		{
			name: "complexe xml",
			xml:  complexeXml,
			path: "feed/entry[0]/updated",
			ret:  []string{"2020-03-27T22:25:57+01:00"},
		},
		{
			name: "complexe xml with attributes",
			xml:  complexeXml,
			path: "feed/entry[1]/link@href",
			ret:  []string{"https://github.com/NeverSinkDev/NeverSink-Filter/releases/tag/7.8.1"},
		},
		{
			name: "complexe xml with depth and attributes",
			xml:  complexeXml,
			path: "feed/link[1]@href",
			ret:  []string{"https://github.com/NeverSinkDev/NeverSink-Filter/releases.atom"},
		},
	} {
		t.Run(d.name, func(t *testing.T) {
			ret := ParseRecursive(d.path, d.xml)
			assert.Equal(t, d.ret, ret)
		})
	}
}
