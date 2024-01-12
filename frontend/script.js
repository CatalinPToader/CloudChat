document.addEventListener("DOMContentLoaded", function() {
    const chatMessages = document.getElementById("chat-messages");
    const messageInput = document.getElementById("message-input");
    const sendButton = document.getElementById("send-button");
    const chatTitle = document.getElementById("chat-title");

    messageInput.disabled = true;

    // Function to send a heartbeat request to the server
    $(document).ready(function() {
        // Function to send a heartbeat request
        function sendHeartbeat() {
            $.ajax({
                url: "http://localhost:3000/heartbeat",  // Replace with your actual server address
                method: "GET",
                success: function(data) {
                    console.log("Heartbeat success:", data);
                },
                error: function(error) {
                    console.error("Heartbeat error:", error);
                }
            });
        }

        // Periodically send heartbeat every 30 seconds (adjust as needed)
        setInterval(sendHeartbeat, 5000);
    });

    function updateChatTitle(channelName) {
        if (channelName) {
            chatTitle.textContent = `Chat Room - ${channelName}`;
            messageInput.disabled = false; // Enable the message input when a channel is selected
        } else {
            chatTitle.textContent = "Chat Room";
            messageInput.disabled = true; // Disable the message input when no channel is selected
        }
    }

    // Get the channel list items
    const channelListItems = document.querySelectorAll("#channels-list li");

    // Add click event listeners to each channel list item
    channelListItems.forEach(function(channelListItem) {
        channelListItem.addEventListener("click", function() {
            const selectedChannel = channelListItem.textContent.trim();

            channelListItems.forEach(function(item) {
                item.classList.remove("selected");
            });

            // Add the "selected" class to the clicked channel list item
            channelListItem.classList.add("selected");

            updateChatTitle(selectedChannel);
        });
    });

    messageInput.addEventListener("keyup", function(event) {
        if (event.key === "Enter" && !event.shiftKey) {
            // Prevent the default behavior of the "Enter" key (newline)
            event.preventDefault();

            // Get the trimmed message from the input
            const message = messageInput.value.trim();

            // If the message is not empty, send it
            if (message !== "") {
                const sender = "You";  // Replace with the actual sender (e.g., username)
                sendMessage(sender, message);

                // Clear the input field
                messageInput.value = "";
            }
        }
    });

    sendButton.addEventListener("click", function() {
        const message = messageInput.value.trim();
        if (message !== "") {
            const sender = "You"; // Replace with the actual sender (e.g., username)
            sendMessage(sender, message);
            messageInput.value = "";
        }
    });

    function appendMessage(sender, message) {
        const messageElement = document.createElement("div");
        messageElement.className = "message";
        messageElement.innerHTML = `<strong>${sender}:</strong> ${message}`;
        chatMessages.appendChild(messageElement);

        // Scroll to the bottom to show the latest message
        chatMessages.scrollTop = chatMessages.scrollHeight;
    }

    function sendMessage(sender, message) {
        const apiUrl = "http://localhost:3000/messages"; // Replace with your actual API endpoint

        // Prepare the data to be sent
        const postData = {
            sender: sender,
            message: message
        };

        // Make an HTTP POST request
        fetch(apiUrl, {
            method: "POST",
            headers: {
                "Content-Type": "application/json"
            },
            body: JSON.stringify(postData)
        })
        .then(response => {
            if (!response.ok) {
                throw new Error("Network response was not ok");
            }
            return response.json();
        })
        .then(data => {
            // Handle the response from the server (if needed)
            console.log("Message sent successfully:", data);
        })
        .catch(error => {
            console.error("Error sending message:", error);
        });

        // Append the message to the chat window immediately
        appendMessage(sender, message);
    }
});
