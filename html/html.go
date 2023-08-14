// Copyright 2015 The Go Authors.
// Copyright 2022 Ryan "rj45" Sanche
// All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package html

import (
	"bytes"
	"crypto/sha256"
	"fmt"
	"html"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/rj45/llbrew/ir"
)

type HTMLWriter struct {
	w             io.WriteCloser
	Func          *ir.Func
	path          string
	prevHash      []byte
	pendingPhases []string
	pendingTitles []string
}

func NewHTMLWriter(path string, f *ir.Func) *HTMLWriter {
	path = strings.Replace(path, "/", string(filepath.Separator), -1)
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
	reportPath := path
	if !filepath.IsAbs(reportPath) {
		pwd, err := os.Getwd()
		if err != nil {
			log.Fatalf("%v", err)
		}
		reportPath = filepath.Join(pwd, path)
	}
	html := HTMLWriter{
		w:    out,
		Func: f,
		path: reportPath,
	}
	html.start()
	return &html
}

// Fatalf reports an error and exits.
func (w *HTMLWriter) Fatalf(msg string, args ...interface{}) {
	log.Fatalf(msg, args...)
}

// Logf calls the (w *HTMLWriter).Func's Logf method passing along a msg and args.
func (w *HTMLWriter) Logf(msg string, args ...interface{}) {
	log.Printf(msg, args...)
}

func (w *HTMLWriter) start() {
	if w == nil {
		return
	}
	w.WriteString("<html>")
	w.WriteString(`<head>
<meta http-equiv="Content-Type" content="text/html;charset=UTF-8">
<style>

body {
    font-size: 14px;
    font-family: Arial, sans-serif;
}

h1 {
    font-size: 18px;
    display: inline-block;
    margin: 0 1em .5em 0;
}

#helplink {
    display: inline-block;
}

#help {
    display: none;
}

.stats {
    font-size: 60%;
}

table {
    border: 1px solid black;
    table-layout: fixed;
    width: 300px;
}

th, td {
    border: 1px solid black;
    overflow: hidden;
    width: 400px;
    vertical-align: top;
    padding: 5px;
}

td > h2 {
    cursor: pointer;
    font-size: 120%;
    margin: 5px 0px 5px 0px;
}

td.collapsed {
    font-size: 12px;
    width: 12px;
    border: 1px solid white;
    padding: 2px;
    cursor: pointer;
    background: #fafafa;
}

td.collapsed div {
    text-align: right;
    transform: rotate(180deg);
    writing-mode: vertical-lr;
    white-space: pre;
}

code, pre, .lines, .ast {
    font-family: Menlo, monospace;
    font-size: 12px;
}

pre {
    -moz-tab-size: 4;
    -o-tab-size:   4;
    tab-size:      4;
}

.allow-x-scroll {
    overflow-x: scroll;
}

.lines {
    float: left;
    overflow: hidden;
    text-align: right;
    margin-top: 7px;
}

.lines div {
    padding-right: 10px;
    color: gray;
}

div.line-number {
    font-size: 12px;
}

.ast {
    white-space: nowrap;
		overflow-x: scroll;

}

td.ssa-prog {
    width: 600px;
    word-wrap: break-word;
}

li {
    list-style-type: none;
}

li.ssa-long-value {
    text-indent: -2em;  /* indent wrapped lines */
}

span.ssa-instr {
	font-weight: bold;
	color: #8250df; // purple
}

span.ssa-instr-copy {
	font-weight: bold;
	color: #bf8700; // yellow
}

span.ssa-instr-call {
	font-weight: bold;
	color: #e85aad; // pink
}

span.ssa-instr-type {
	color: #c5cacf; // grey
}

span.ssa-value {
	width: 10em;
}

span.ssa-value-id {
	color: #0969da; // blue
}

span.ssa-value-reg {
	color: #cf222e; // red
}

span.ssa-value-const {
	color: #8c959f; // grey
}

span.ssa-value-const-num {
	color: #2da44e; // green
}

li.ssa-value-list {
    display: inline;
}

li.ssa-start-block {
    padding: 0;
    margin: 0;
}

li.ssa-end-block {
    padding: 0;
    margin: 0;
}

ul.ssa-print-func {
    padding-left: 0;
}

li.ssa-start-block button {
    padding: 0 1em;
    margin: 0;
    border: none;
    display: inline;
    font-size: 14px;
    float: right;
}

button:hover {
    background-color: #eee;
    cursor: pointer;
}

dl.ssa-gen {
    padding-left: 0;
}

dt.ssa-prog-src {
    padding: 0;
    margin: 0;
    float: left;
    width: 4em;
}

dd.ssa-prog {
    padding: 0;
    margin-right: 0;
    margin-left: 4em;
}

.dead-value {
    color: gray;
}

.dead-block {
    opacity: 0.5;
}

.depcycle {
    font-style: italic;
}

.line-number {
    font-size: 11px;
}

.no-line-number {
    font-size: 11px;
    color: gray;
}

.zoom {
	position: absolute;
	float: left;
	white-space: nowrap;
	background-color: #eee;
}

.zoom a:link, .zoom a:visited  {
    text-decoration: none;
    color: blue;
    font-size: 16px;
    padding: 4px 2px;
}

svg {
    cursor: default;
    outline: 1px solid #eee;
    width: 100%;
}

body.darkmode {
    background-color: rgb(21, 21, 21);
    color: rgb(230, 255, 255);
    opacity: 100%;
}

td.darkmode {
    background-color: rgb(21, 21, 21);
    border: 1px solid gray;
}

body.darkmode table, th {
    border: 1px solid gray;
}

body.darkmode text {
    fill: white;
}

body.darkmode svg polygon:first-child {
    fill: rgb(21, 21, 21);
}

.highlight-aquamarine     { background-color: aquamarine; color: black; }
.highlight-coral          { background-color: coral; color: black; }
.highlight-lightpink      { background-color: lightpink; color: black; }
.highlight-lightsteelblue { background-color: lightsteelblue; color: black; }
.highlight-palegreen      { background-color: palegreen; color: black; }
.highlight-skyblue        { background-color: skyblue; color: black; }
.highlight-lightgray      { background-color: lightgray; color: black; }
.highlight-yellow         { background-color: yellow; color: black; }
.highlight-lime           { background-color: lime; color: black; }
.highlight-khaki          { background-color: khaki; color: black; }
.highlight-aqua           { background-color: aqua; color: black; }
.highlight-salmon         { background-color: salmon; color: black; }

/* Ensure all dead values/blocks continue to have gray font color in dark mode with highlights */
.dead-value span.highlight-aquamarine,
.dead-block.highlight-aquamarine,
.dead-value span.highlight-coral,
.dead-block.highlight-coral,
.dead-value span.highlight-lightpink,
.dead-block.highlight-lightpink,
.dead-value span.highlight-lightsteelblue,
.dead-block.highlight-lightsteelblue,
.dead-value span.highlight-palegreen,
.dead-block.highlight-palegreen,
.dead-value span.highlight-skyblue,
.dead-block.highlight-skyblue,
.dead-value span.highlight-lightgray,
.dead-block.highlight-lightgray,
.dead-value span.highlight-yellow,
.dead-block.highlight-yellow,
.dead-value span.highlight-lime,
.dead-block.highlight-lime,
.dead-value span.highlight-khaki,
.dead-block.highlight-khaki,
.dead-value span.highlight-aqua,
.dead-block.highlight-aqua,
.dead-value span.highlight-salmon,
.dead-block.highlight-salmon {
    color: gray;
}

.outline-blue           { outline: #2893ff solid 2px; }
.outline-red            { outline: red solid 2px; }
.outline-blueviolet     { outline: blueviolet solid 2px; }
.outline-darkolivegreen { outline: darkolivegreen solid 2px; }
.outline-fuchsia        { outline: fuchsia solid 2px; }
.outline-sienna         { outline: sienna solid 2px; }
.outline-gold           { outline: gold solid 2px; }
.outline-orangered      { outline: orangered solid 2px; }
.outline-teal           { outline: teal solid 2px; }
.outline-maroon         { outline: maroon solid 2px; }
.outline-black          { outline: black solid 2px; }

ellipse.outline-blue           { stroke-width: 2px; stroke: #2893ff; }
ellipse.outline-red            { stroke-width: 2px; stroke: red; }
ellipse.outline-blueviolet     { stroke-width: 2px; stroke: blueviolet; }
ellipse.outline-darkolivegreen { stroke-width: 2px; stroke: darkolivegreen; }
ellipse.outline-fuchsia        { stroke-width: 2px; stroke: fuchsia; }
ellipse.outline-sienna         { stroke-width: 2px; stroke: sienna; }
ellipse.outline-gold           { stroke-width: 2px; stroke: gold; }
ellipse.outline-orangered      { stroke-width: 2px; stroke: orangered; }
ellipse.outline-teal           { stroke-width: 2px; stroke: teal; }
ellipse.outline-maroon         { stroke-width: 2px; stroke: maroon; }
ellipse.outline-black          { stroke-width: 2px; stroke: black; }

/* Capture alternative for outline-black and ellipse.outline-black when in dark mode */
body.darkmode .outline-black        { outline: gray solid 2px; }
body.darkmode ellipse.outline-black { outline: gray solid 2px; }

</style>

<script type="text/javascript">

// Contains phase names which are expanded by default. Other columns are collapsed.
let expandedDefault = [
    "start",
    "deadcode",
    "opt",
    "lower",
    "late-deadcode",
    "regalloc",
    "genssa",
];
if (history.state === null) {
    history.pushState({expandedDefault}, "", location.href);
}

// ordered list of all available highlight colors
var highlights = [
    "highlight-aquamarine",
    "highlight-coral",
    "highlight-lightpink",
    "highlight-lightsteelblue",
    "highlight-palegreen",
    "highlight-skyblue",
    "highlight-lightgray",
    "highlight-yellow",
    "highlight-lime",
    "highlight-khaki",
    "highlight-aqua",
    "highlight-salmon"
];

// state: which value is highlighted this color?
var highlighted = {};
for (var i = 0; i < highlights.length; i++) {
    highlighted[highlights[i]] = "";
}

// ordered list of all available outline colors
var outlines = [
    "outline-blue",
    "outline-red",
    "outline-blueviolet",
    "outline-darkolivegreen",
    "outline-fuchsia",
    "outline-sienna",
    "outline-gold",
    "outline-orangered",
    "outline-teal",
    "outline-maroon",
    "outline-black"
];

// state: which value is outlined this color?
var outlined = {};
for (var i = 0; i < outlines.length; i++) {
    outlined[outlines[i]] = "";
}

window.onload = function() {
    if (history.state !== null) {
        expandedDefault = history.state.expandedDefault;
    }
    if (window.matchMedia && window.matchMedia("(prefers-color-scheme: dark)").matches) {
        toggleDarkMode();
        document.getElementById("dark-mode-button").checked = true;
    }

    var ssaElemClicked = function(elem, event, selections, selected) {
        event.stopPropagation();

        // find all values with the same name
        var c = elem.classList.item(0);
        var x = document.getElementsByClassName(c);

        // if selected, remove selections from all of them
        // otherwise, attempt to add

        var remove = "";
        for (var i = 0; i < selections.length; i++) {
            var color = selections[i];
            if (selected[color] == c) {
                remove = color;
                break;
            }
        }

        if (remove != "") {
            for (var i = 0; i < x.length; i++) {
                x[i].classList.remove(remove);
            }
            selected[remove] = "";
            return;
        }

        // we're adding a selection
        // find first available color
        var avail = "";
        for (var i = 0; i < selections.length; i++) {
            var color = selections[i];
            if (selected[color] == "") {
                avail = color;
                break;
            }
        }
        if (avail == "") {
            alert("out of selection colors; go add more");
            return;
        }

        // set that as the selection
        for (var i = 0; i < x.length; i++) {
            x[i].classList.add(avail);
        }
        selected[avail] = c;
    };

    var ssaValueClicked = function(event) {
        ssaElemClicked(this, event, highlights, highlighted);
    };

		var ssaValueRegClicked = function(event) {
				ssaElemClicked(this, event, highlights, highlighted);
		};

    var ssaBlockClicked = function(event) {
        ssaElemClicked(this, event, outlines, outlined);
    };

    var ssavalues = document.getElementsByClassName("ssa-value");
    for (var i = 0; i < ssavalues.length; i++) {
        ssavalues[i].addEventListener('click', ssaValueClicked);
    }

    var ssavalueregs = document.getElementsByClassName("ssa-value-reg");
    for (var i = 0; i < ssavalueregs.length; i++) {
				ssavalueregs[i].addEventListener('click', ssaValueRegClicked);
    }

    var ssalongvalues = document.getElementsByClassName("ssa-long-value");
    for (var i = 0; i < ssalongvalues.length; i++) {
        // don't attach listeners to li nodes, just the spans they contain
        if (ssalongvalues[i].nodeName == "SPAN") {
            ssalongvalues[i].addEventListener('click', ssaValueClicked);
        }
    }

    var ssablocks = document.getElementsByClassName("ssa-block");
    for (var i = 0; i < ssablocks.length; i++) {
        ssablocks[i].addEventListener('click', ssaBlockClicked);
    }

    var lines = document.getElementsByClassName("line-number");
    for (var i = 0; i < lines.length; i++) {
        lines[i].addEventListener('click', ssaValueClicked);
    }


    function toggler(phase) {
        return function() {
            toggle_cell(phase+'-col');
            toggle_cell(phase+'-exp');
            const i = expandedDefault.indexOf(phase);
            if (i !== -1) {
                expandedDefault.splice(i, 1);
            } else {
                expandedDefault.push(phase);
            }
            history.pushState({expandedDefault}, "", location.href);
        };
    }

    function toggle_cell(id) {
        var e = document.getElementById(id);
        if (e.style.display == 'table-cell') {
            e.style.display = 'none';
        } else {
            e.style.display = 'table-cell';
        }
    }

    // Go through all columns and collapse needed phases.
    const td = document.getElementsByTagName("td");
    for (let i = 0; i < td.length; i++) {
        const id = td[i].id;
        const phase = id.substr(0, id.length-4);
        let show = expandedDefault.indexOf(phase) !== -1

        // If show == false, check to see if this is a combined column (multiple phases).
        // If combined, check each of the phases to see if they are in our expandedDefaults.
        // If any are found, that entire combined column gets shown.
        if (!show) {
            const combined = phase.split('--+--');
            const len = combined.length;
            if (len > 1) {
                for (let i = 0; i < len; i++) {
                    const num = expandedDefault.indexOf(combined[i]);
                    if (num !== -1) {
                        expandedDefault.splice(num, 1);
                        if (expandedDefault.indexOf(phase) === -1) {
                            expandedDefault.push(phase);
                            show = true;
                        }
                    }
                }
            }
        }
        if (id.endsWith("-exp")) {
            const h2Els = td[i].getElementsByTagName("h2");
            const len = h2Els.length;
            if (len > 0) {
                for (let i = 0; i < len; i++) {
                    h2Els[i].addEventListener('click', toggler(phase));
                }
            }
        } else {
            td[i].addEventListener('click', toggler(phase));
        }
        if (id.endsWith("-col") && show || id.endsWith("-exp") && !show) {
            td[i].style.display = 'none';
            continue;
        }
        td[i].style.display = 'table-cell';
    }

    // find all svg block nodes, add their block classes
    var nodes = document.querySelectorAll('*[id^="graph_node_"]');
    for (var i = 0; i < nodes.length; i++) {
    	var node = nodes[i];
    	var name = node.id.toString();
    	var block = name.substring(name.lastIndexOf("_")+1);
    	node.classList.remove("node");
    	node.classList.add(block);
        node.addEventListener('click', ssaBlockClicked);
        var ellipse = node.getElementsByTagName('ellipse')[0];
        ellipse.classList.add(block);
        ellipse.addEventListener('click', ssaBlockClicked);
    }

    // make big graphs smaller
    var targetScale = 0.5;
    var nodes = document.querySelectorAll('*[id^="svg_graph_"]');
    // TODO: Implement smarter auto-zoom using the viewBox attribute
    // and in case of big graphs set the width and height of the svg graph to
    // maximum allowed.
    for (var i = 0; i < nodes.length; i++) {
    	var node = nodes[i];
    	var name = node.id.toString();
    	var phase = name.substring(name.lastIndexOf("_")+1);
    	var gNode = document.getElementById("g_graph_"+phase);
    	var scale = gNode.transform.baseVal.getItem(0).matrix.a;
    	if (scale > targetScale) {
    		node.width.baseVal.value *= targetScale / scale;
    		node.height.baseVal.value *= targetScale / scale;
    	}
    }
};

function toggle_visibility(id) {
    var e = document.getElementById(id);
    if (e.style.display == 'block') {
        e.style.display = 'none';
    } else {
        e.style.display = 'block';
    }
}

function hideBlock(el) {
    var es = el.parentNode.parentNode.getElementsByClassName("ssa-value-list");
    if (es.length===0)
        return;
    var e = es[0];
    if (e.style.display === 'block' || e.style.display === '') {
        e.style.display = 'none';
        el.innerHTML = '+';
    } else {
        e.style.display = 'block';
        el.innerHTML = '-';
    }
}

// TODO: scale the graph with the viewBox attribute.
function graphReduce(id) {
    var node = document.getElementById(id);
    if (node) {
    		node.width.baseVal.value *= 0.9;
    		node.height.baseVal.value *= 0.9;
    }
    return false;
}

function graphEnlarge(id) {
    var node = document.getElementById(id);
    if (node) {
    		node.width.baseVal.value *= 1.1;
    		node.height.baseVal.value *= 1.1;
    }
    return false;
}

function makeDraggable(event) {
    var svg = event.target;
    if (window.PointerEvent) {
        svg.addEventListener('pointerdown', startDrag);
        svg.addEventListener('pointermove', drag);
        svg.addEventListener('pointerup', endDrag);
        svg.addEventListener('pointerleave', endDrag);
    } else {
        svg.addEventListener('mousedown', startDrag);
        svg.addEventListener('mousemove', drag);
        svg.addEventListener('mouseup', endDrag);
        svg.addEventListener('mouseleave', endDrag);
    }

    var point = svg.createSVGPoint();
    var isPointerDown = false;
    var pointerOrigin;
    var viewBox = svg.viewBox.baseVal;

    function getPointFromEvent (event) {
        point.x = event.clientX;
        point.y = event.clientY;

        // We get the current transformation matrix of the SVG and we inverse it
        var invertedSVGMatrix = svg.getScreenCTM().inverse();
        return point.matrixTransform(invertedSVGMatrix);
    }

    function startDrag(event) {
        isPointerDown = true;
        pointerOrigin = getPointFromEvent(event);
    }

    function drag(event) {
        if (!isPointerDown) {
            return;
        }
        event.preventDefault();

        var pointerPosition = getPointFromEvent(event);
        viewBox.x -= (pointerPosition.x - pointerOrigin.x);
        viewBox.y -= (pointerPosition.y - pointerOrigin.y);
    }

    function endDrag(event) {
        isPointerDown = false;
    }
}

function toggleDarkMode() {
    document.body.classList.toggle('darkmode');

    // Collect all of the "collapsed" elements and apply dark mode on each collapsed column
    const collapsedEls = document.getElementsByClassName('collapsed');
    const len = collapsedEls.length;

    for (let i = 0; i < len; i++) {
        collapsedEls[i].classList.toggle('darkmode');
    }

    // Collect and spread the appropriate elements from all of the svgs on the page into one array
    const svgParts = [
        ...document.querySelectorAll('path'),
        ...document.querySelectorAll('ellipse'),
        ...document.querySelectorAll('polygon'),
    ];

    // Iterate over the svgParts specifically looking for white and black fill/stroke to be toggled.
    // The verbose conditional is intentional here so that we do not mutate any svg path, ellipse, or polygon that is of any color other than white or black.
    svgParts.forEach(el => {
        if (el.attributes.stroke.value === 'white') {
            el.attributes.stroke.value = 'black';
        } else if (el.attributes.stroke.value === 'black') {
            el.attributes.stroke.value = 'white';
        }
        if (el.attributes.fill.value === 'white') {
            el.attributes.fill.value = 'black';
        } else if (el.attributes.fill.value === 'black') {
            el.attributes.fill.value = 'white';
        }
    });
}

</script>

</head>`)
	w.WriteString("<body>")
	w.WriteString("<h1>")
	w.WriteString(html.EscapeString(w.Func.Name))
	w.WriteString("</h1>")
	w.WriteString(`
<a href="#" onclick="toggle_visibility('help');return false;" id="helplink">help</a>
<div id="help">

<p>
Click on a value or block to toggle highlighting of that value/block
and its uses.  (Values and blocks are highlighted by ID, and IDs of
dead items may be reused, so not all highlights necessarily correspond
to the clicked item.)
</p>

<p>
Faded out values and blocks are dead code that has not been eliminated.
</p>

<p>
Values printed in italics have a dependency cycle.
</p>

<p>
<b>CFG</b>: Dashed edge is for unlikely branches. Blue color is for backward edges.
Edge with a dot means that this edge follows the order in which blocks were laidout.
</p>

</div>
<label for="dark-mode-button" style="margin-left: 15px; cursor: pointer;">darkmode</label>
<input type="checkbox" onclick="toggleDarkMode();" id="dark-mode-button" style="cursor: pointer" />
`)
	w.WriteString("<table>")
	w.WriteString("<tr>")
}

func (w *HTMLWriter) Close() {
	if w == nil {
		return
	}
	io.WriteString(w.w, "</tr>")
	io.WriteString(w.w, "</table>")
	io.WriteString(w.w, "</body>")
	io.WriteString(w.w, "</html>")
	w.w.Close()
	fmt.Printf("dumped SSA to %v\n", w.path)
}

func hashFunc(f *ir.Func) []byte {
	h := sha256.New()
	f.Emit(h, IRDecorator{})
	return h.Sum(nil)
}

// WritePhase writes f in a column headed by title.
// phase is used for collapsing columns and should be unique across the table.
func (w *HTMLWriter) WritePhase(phase, title string) {
	if w == nil {
		return // avoid generating HTML just to discard it
	}
	hash := hashFunc(w.Func)
	w.pendingPhases = append(w.pendingPhases, phase)
	w.pendingTitles = append(w.pendingTitles, title)
	if !bytes.Equal(hash, w.prevHash) {
		w.flushPhases()
	}
	w.prevHash = hash
}

// flushPhases collects any pending phases and titles, writes them to the html, and resets the pending slices.
func (w *HTMLWriter) flushPhases() {
	phaseLen := len(w.pendingPhases)
	if phaseLen == 0 {
		return
	}
	phases := strings.Join(w.pendingPhases, "  +  ")
	w.WriteMultiTitleColumn(
		phases,
		w.pendingTitles,
		fmt.Sprintf("hash-%x", w.prevHash),
		FuncHTML(w.Func, w.pendingPhases[phaseLen-1]),
	)
	w.pendingPhases = w.pendingPhases[:0]
	w.pendingTitles = w.pendingTitles[:0]
}

func FuncHTML(f *ir.Func, phase string) string {
	buf := new(bytes.Buffer)
	fmt.Fprint(buf, "<code>")
	f.Emit(buf, IRDecorator{})
	fmt.Fprint(buf, "</code>")
	return buf.String()
}

// FuncLines contains source code for a function to be displayed
// in sources column.
type FuncLines struct {
	Filename    string
	StartLineno uint
	Lines       []string
}

// ByTopo sorts topologically: target function is on top,
// followed by inlined functions sorted by filename and line numbers.
type ByTopo []FuncLines

func (x ByTopo) Len() int      { return len(x) }
func (x ByTopo) Swap(i, j int) { x[i], x[j] = x[j], x[i] }
func (x ByTopo) Less(i, j int) bool {
	a := x[i]
	b := x[j]
	if a.Filename == b.Filename {
		return a.StartLineno < b.StartLineno
	}
	return a.Filename < b.Filename
}

// WriteSources writes lines as source code in a column headed by title.
// phase is used for collapsing columns and should be unique across the table.
func (w *HTMLWriter) WriteSources(phase string, fn string, lines []string, startline int) {
	if fn == "" || len(lines) == 0 {
		return
	}

	all := []*FuncLines{
		{
			Filename:    fn,
			StartLineno: uint(startline),
			Lines:       lines,
		},
	}

	if w == nil {
		return // avoid generating HTML just to discard it
	}
	var buf bytes.Buffer
	fmt.Fprint(&buf, "<div class=\"lines\" style=\"width: 8%\">")
	filename := ""
	for _, fl := range all {
		fmt.Fprint(&buf, "<div>&nbsp;</div>")
		if filename != fl.Filename {
			fmt.Fprint(&buf, "<div>&nbsp;</div>")
			filename = fl.Filename
		}
		for i := range fl.Lines {
			ln := int(fl.StartLineno) + i
			fmt.Fprintf(&buf, "<div class=\"l%v line-number\">%v</div>", ln, ln)
		}
	}
	fmt.Fprint(&buf, "</div><div class=\"allow-x-scroll\" style=\"width: 92%\"><pre>")
	filename = ""
	for _, fl := range all {
		fmt.Fprint(&buf, "<div>&nbsp;</div>")
		if filename != fl.Filename {
			fmt.Fprintf(&buf, "<div><strong>%v</strong></div>", fl.Filename)
			filename = fl.Filename
		}
		for i, line := range fl.Lines {
			ln := int(fl.StartLineno) + i
			var escaped string
			if strings.TrimSpace(line) == "" {
				escaped = "&nbsp;"
			} else {
				escaped = html.EscapeString(line)
			}
			fmt.Fprintf(&buf, "<div class=\"l%v line-number\">%v</div>", ln, escaped)
		}
	}
	fmt.Fprint(&buf, "</pre></div>")
	w.WriteColumn(phase, phase, "allow-x-scroll", buf.String())
}

func (w *HTMLWriter) WriteAST(phase string, buf *bytes.Buffer) {
	if w == nil {
		return // avoid generating HTML just to discard it
	}
	lines := strings.Split(buf.String(), "\n")
	var out bytes.Buffer

	fmt.Fprint(&out, "<div>")
	for _, l := range lines {
		l = strings.TrimSpace(l)
		var escaped string
		var lineNo string
		if l == "" {
			escaped = "&nbsp;"
		} else {
			if strings.HasPrefix(l, "buildssa") {
				escaped = fmt.Sprintf("<b>%v</b>", l)
			} else {
				// Parse the line number from the format file:line:col.
				// See the implementation in ir/fmt.go:dumpNodeHeader.
				sl := strings.Split(l, ":")
				if len(sl) >= 3 {
					if _, err := strconv.Atoi(sl[len(sl)-2]); err == nil {
						lineNo = sl[len(sl)-2]
					}
				}
				escaped = html.EscapeString(l)
			}
		}
		if lineNo != "" {
			fmt.Fprintf(&out, "<div class=\"l%v line-number ast\">%v</div>", lineNo, escaped)
		} else {
			fmt.Fprintf(&out, "<div class=\"ast\">%v</div>", escaped)
		}
	}
	fmt.Fprint(&out, "</div>")
	w.WriteColumn(phase, phase, "allow-x-scroll", out.String())
}

// WriteColumn writes raw HTML in a column headed by title.
// It is intended for pre- and post-compilation log output.
func (w *HTMLWriter) WriteColumn(phase, title, class, html string) {
	w.WriteMultiTitleColumn(phase, []string{title}, class, html)
}

func (w *HTMLWriter) WriteMultiTitleColumn(phase string, titles []string, class, html string) {
	if w == nil {
		return
	}
	id := strings.Replace(phase, " ", "-", -1)
	// collapsed column
	w.Printf("<td id=\"%v-col\" class=\"collapsed\"><div>%v</div></td>", id, phase)

	if class == "" {
		w.Printf("<td id=\"%v-exp\">", id)
	} else {
		w.Printf("<td id=\"%v-exp\" class=\"%v\">", id, class)
	}
	for _, title := range titles {
		w.WriteString("<h2>" + title + "</h2>")
	}
	w.WriteString(html)
	w.WriteString("</td>\n")
}

func (w *HTMLWriter) Printf(msg string, v ...interface{}) {
	if _, err := fmt.Fprintf(w.w, msg, v...); err != nil {
		w.Fatalf("%v", err)
	}
}

func (w *HTMLWriter) WriteString(s string) {
	if _, err := io.WriteString(w.w, s); err != nil {
		w.Fatalf("%v", err)
	}
}
