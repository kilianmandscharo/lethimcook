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

export function insertStringInInputAtCursor(
    el: HTMLTextAreaElement,
    s: string,
) {
    const cursorPos = el.selectionStart;
    const newCursorPos = cursorPos + s.length;
    el.value =
        el.value.slice(0, cursorPos) + s + el.value.slice(el.selectionEnd);
    el.selectionStart = el.selectionEnd = newCursorPos;
}

export function extractContentBetweenBrackets<
    T extends {
        selectionStart: number;
        value: string;
    },
>(
    target: T,
): { substitutionEnd: number; substitutionStart: number; query: string } | null {
    if (target.value.length === 0) {
        return null;
    }

    const cursorPos = target.selectionStart;
    const lastChar = target.value[cursorPos - 1];
    if (lastChar !== "]") {
        return null;
    }

    let bracketContentReversed = "";
    let i = cursorPos - 2;
    let closingPatternFound = false;
    while (i > 0) {
        if (target.value[i] === "!" && target.value[i - 1] === "[") {
            closingPatternFound = true;
            break;
        }
        bracketContentReversed += target.value[i];
        i--;
    }

    if (!closingPatternFound) {
        return null;
    }

    return {
        query: bracketContentReversed.split("").reverse().join(""),
        substitutionStart: i - 1,
        substitutionEnd: cursorPos,
    };
}

export function fetchLink(target: HTMLTextAreaElement) {
    const extractionData = extractContentBetweenBrackets(target);
    if (!extractionData) {
        return;
    }

    const params = new URLSearchParams({
        query: extractionData.query,
    });

    fetch(`${window.location.origin}/recipe/link?` + params.toString())
        .then((res) => res.text())
        .then((data) => {
            SelectDialog.state.target = target;
            SelectDialog.state.cursorPos = extractionData.substitutionEnd;
            SelectDialog.state.substitutionStart =
                extractionData.substitutionStart;
            try {
                const recipe: { title: string; id: number } = JSON.parse(data);
                SelectDialog.injectLinkIntoTextarea(recipe.title, recipe.id);
            } catch {
                SelectDialog.open(data);
            }
        })
        .catch(console.error);
}
