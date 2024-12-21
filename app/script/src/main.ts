import {
    attachTextAreaEventListeners,
    handleCopyUrlToClipboardButton,
} from "./globals";
import { LocalStorageUtil } from "./localStorage";
import { attachMutationObserverToNotificationContainer } from "./notification";
import { SelectDialog } from "./selectDialog";

function main() {
    attachMutationObserverToNotificationContainer();
    (window as any).LocalStorageUtil = LocalStorageUtil;
    (window as any).SelectDialog = SelectDialog;
    (window as any).attachTextAreaEventListeners = attachTextAreaEventListeners;
    (window as any).handleCopyUrlToClipboardButton =
        handleCopyUrlToClipboardButton;
}

main();
