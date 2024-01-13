// main.js

$(document).ready(function () {
    // Send a request to fetch channels when the document is loaded
    fetchChannels();

    clearUserList();

    // Send a heartbeat every second
    setInterval(sendHeartbeat, 1000);
});
