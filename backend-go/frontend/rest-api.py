from flask import Flask, request, jsonify
from flask_cors import CORS  # Import CORS from flask_cors

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

@app.route('/messages', methods=['POST'])
def send_message():
    try:
        data = request.get_json()

        # Assuming the data contains "sender" and "message" fields
        sender = data.get('sender', '')
        message = data.get('message', '')

        # Print sender and message to the console (for debugging)
        print(f"Received message from {sender}: {message}")

        # Process the message (e.g., save it to a database)

        return jsonify({'status': 'success', 'message': 'Message received successfully'}), 200
    except Exception as e:
        return jsonify({'status': 'error', 'message': str(e)}), 500

@app.route('/heartbeat', methods=['GET'])
def heartbeat():
    # Optionally, you can perform some logging or additional checks here
    print('Received heartbeat request')
    
    # Return a simple response to acknowledge the heartbeat
    return jsonify({'status': 'success', 'message': 'Heartbeat received'})


if __name__ == '__main__':
    # Set the port to your desired value, e.g., 8080
    app.run(debug=True, port=3000)
