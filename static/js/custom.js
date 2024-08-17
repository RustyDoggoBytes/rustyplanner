function clearAndFocus(event, elementId) {
    const $input = document.getElementById(elementId);
    document.addEventListener('DOMContentLoaded', function () {
        $input.focus();
    });

    if (event.detail.successful) {
        $input.value = '';
        // Delay focus for iOS devices
        setTimeout(function () {
            $input.focus();
            // Scroll to the input if needed
            $input.scrollIntoView({behavior: "smooth", block: "center"});
        }, 100);
    }
}
