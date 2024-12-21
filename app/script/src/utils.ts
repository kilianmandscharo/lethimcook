import { SelectDialog } from "./selectDialog";

export function replaceScripts(node: Element) {
    if (node.tagName === "SCRIPT") {
        const script = document.createElement("script");
        script.text = node.innerHTML;
        for (let i = 0; i < node.attributes.length; i++) {
            script.setAttribute(
                node.attributes[i].name,
                node.attributes[i].value,
            );
        }
        node.parentNode?.replaceChild(script, node);
    }
    for (let i = 0; i < node.children.length; i++) {
        const el = node.children.item(i);
        if (el) {
            replaceScripts(el);
        }
    }
}

export function getLastListNumberBeforeCursor(el: HTMLTextAreaElement) {
    const re = /(\n(\d+)|^(\d+))\./g;
    const matches = [...el.value.slice(0, el.selectionStart).matchAll(re)];
    return matches.length > 0 ? matches[matches.length - 1][1] : null;
}

export function insertStringInInputAtCursor(el: HTMLTextAreaElement, s: string) {
    const cursorPos = el.selectionStart;
    const newCursorPos = cursorPos + s.length;
    el.value =
        el.value.slice(0, cursorPos) + s + el.value.slice(el.selectionEnd);
    el.selectionStart = el.selectionEnd = newCursorPos;
}

export function fetchLink(target: HTMLTextAreaElement) {
    const cursorPos = target.selectionStart;
    if (target.value.length === 0) {
        return;
    }
    const lastChar = target.value[cursorPos - 1];
    if (lastChar !== "]") {
        return;
    }
    let bracketContentReversed = "";
    let i = cursorPos - 2;
    while (i-- > 0) {
        if (target.value[i] === "[") {
            break;
        }
        bracketContentReversed += target.value[i];
    }
    const params = new URLSearchParams({
        query: bracketContentReversed.split("").reverse().join(""),
    });
    fetch(`${window.location.origin}/recipe/link?` + params.toString())
        .then((res) => res.text())
        .then((data) => {
            SelectDialog.state.target = target;
            SelectDialog.state.cursorPos = cursorPos;
            SelectDialog.state.substitutionStart = i;
            try {
                const recipe = JSON.parse(data);
                SelectDialog.injectLinkIntoTextarea(recipe.title, recipe.id);
            } catch {
                SelectDialog.open(data);
            }
        })
        .catch(console.error);
}
