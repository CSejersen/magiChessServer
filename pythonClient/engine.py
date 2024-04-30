import asyncio
import websockets
import time

async def connect_to_server():
    # Connect to your WebSocket server
    uri = "ws://localhost:3000/ws"
    async with websockets.connect(uri) as websocket:

        print("Sending identification msg to server")
        await websocket.send("engine")
        ack = await websocket.recv()
        print(f"from server {ack}")

        while(True):
            print("Waiting for task from server")
            task = await websocket.recv()
            print("got task: ", task)

            if (task[0] == "l"):
                await websocket.send("le2")
            if (task[0] == "m"):
                await websocket.send("me2e4")

            time.sleep(4)


# Run the asyncio event loop
asyncio.run(connect_to_server())

