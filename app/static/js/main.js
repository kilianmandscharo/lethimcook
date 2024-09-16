function main() {
    addErrorResponseListener();
    addMessageListener();
}

function addErrorResponseListener() {
    document.body.addEventListener("htmx:responseError", function(event) {
        const message = event.detail.xhr.response;
        notify(message, "error");
    });
}

function addMessageListener() {
    document.body.addEventListener("message", function(event) {
        const message = JSON.parse(event.detail.value);
        if (message.isError) {
            notify(message.value, "error");
        } else {
            notify(message.value, "success");
        }
    });
}

function notify(message, type) {
    const notificationContainer = document.getElementById(
        "notification-container",
    );
    if (!notificationContainer) {
        return;
    }

    const notifications = document.querySelectorAll(".notification");
    const index = notifications.length;

    const newNotification = document.createElement("div");
    newNotification.id = `notification-${index}`;
    newNotification.className =
        type === "success"
            ? "notification notification-success"
            : "notification notification-danger";

    const newNotificationIcon = document.createElement("div");
    newNotificationIcon.className =
        type === "success"
            ? "fa-solid fa-circle-info"
            : "fa-regular fa-circle-exclamation";

    const newNotificationText = document.createElement("p");
    newNotificationText.className = "notification-text";
    newNotificationText.innerText = message;

    newNotification.appendChild(newNotificationIcon);
    newNotification.appendChild(newNotificationText);

    notificationContainer.appendChild(newNotification);

    setTimeout(() => {
        notificationContainer.removeChild(newNotification);
    }, 3000);
}

main();
