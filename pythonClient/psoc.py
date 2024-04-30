import asyncio
import websockets
import time

async def connect_to_server():
    # Connect to your WebSocket server
    uri = "ws://localhost:3000/ws"
    async with websockets.connect(uri) as websocket:

        print("Sending identification msg to server")
        await websocket.send("psoc")
        ack = await websocket.recv()
        print(f"from server {ack}")

        while(True):
            print("Waiting for task from server")
            task = await websocket.recv()
            print("got task: ", task)
            await websocket.send("y")


# Run the asyncio event loop
asyncio.run(connect_to_server())



