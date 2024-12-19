
CodeMirror.defineSimpleMode("gfxscript", {
    start: [
        { regex: /\`.*?\`/, token: "string" },
        { regex: /\b(use)\b/, token: "keyword" },
        { regex: /\b(add|average|color-burn|darken|difference|divide|erase|exclusion|hard-light|lighten|linear-burn|multiply|negation|normal|overlay|pin-light|screen|soft-light|subtract)\b/, token: "atom" },
        { regex: /\b(name|width|height|color|crop|filter|resize)\b(?=\s*=)/, token: "attribute" },
        { regex: /\b-?\d+(\.\d+)?\b/, token: "number" },
        { regex: /#.*$/, token: "comment" },
        { regex: /\[\b(VARS|FILTERS|COMPOSITION|LAYERS)\b\]/, token: "def" },
    ],
    meta: {
        lineComment: "#"
    }
});
