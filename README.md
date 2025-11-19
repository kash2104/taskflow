# Distributed Code Executor

A scalable Master–Worker based distributed system for executing untrusted user code in isolated sandboxes.

## Features

**Master-Worker Architecture**

- Master exposes HTTP APIs for code submission & result retrieval.
- Workers pull tasks via RPC and execute them securely.

**Task Queue + Pending Queue**

- _Task Queue_: Stores tasks waiting for workers.
- _Pending Queue_: Tracks tasks currently being executed.
- Ensures retries even if a worker crashes.

**Secure Sandboxed Execution**

- Each task runs in a separate isolated environment.

## Architecture

<img width="1710" height="591" alt="Image" src="https://github.com/user-attachments/assets/0b620805-2656-4b8c-91b9-e937725305de" />

---

<img width="1374" height="763" alt="Image" src="https://github.com/user-attachments/assets/f92bf085-9195-4a54-ab05-d2f2beb0ef47" />
