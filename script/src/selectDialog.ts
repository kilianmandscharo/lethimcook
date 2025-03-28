import { replaceScripts } from "./utils";

type State = {
    cursorPos?: number;
    target?: HTMLTextAreaElement;
    substitutionStart?: number;
};

export class SelectDialog {
    static state: State = {};

    static injectLinkIntoTextarea(title: string, id: number) {
        const { cursorPos, target, substitutionStart } = SelectDialog.state;
        if (
            cursorPos === undefined ||
            target === undefined ||
            substitutionStart === undefined
        ) {
            return;
        }
        const link = `[${title}](${window.location.origin}/recipe/${id})`;
        const newCursorPos = cursorPos + link.length;
        target.value =
            target.value.slice(0, substitutionStart) +
            link +
            target.value.slice(SelectDialog.state.cursorPos);
        target.selectionStart = target.selectionEnd =
            newCursorPos - (cursorPos - substitutionStart);
        SelectDialog.state = {};
    }

    static close() {
        document.getElementById("select-dialog-container")?.remove();
    }

    static open(data: string) {
        const container = document.createElement("div");
        container.id = "select-dialog-container"
        container.innerHTML = data;
        replaceScripts(container);
        document.body.appendChild(container);
    }
}
