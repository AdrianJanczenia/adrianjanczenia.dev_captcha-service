# Adrian Janczenia - Captcha Service

The **Captcha Service** acts as a sophisticated security gateway within the portfolio microservices ecosystem. It provides a multi-layered defense mechanism designed to distinguish human users from automated bots through a hybrid approach combining cryptographic puzzles and visual challenges.

## Service Role

This service is the "shield" of the ecosystem, ensuring that resource-heavy operations (like sending emails or accessing private data) are protected. Its primary responsibilities include:

- **Proof of Work (PoW) Generation**: Issuing cryptographically signed "seeds" that require computational effort from the client before a captcha is even presented.
- **Visual Captcha Orchestration**: Generating unique, time-bound visual puzzles and managing their state (answers, remaining tries, and expiration).
- **Double-Spend Prevention**: Ensuring that a single PoW solution or a solved captcha cannot be reused, effectively neutralizing replay attacks.
- **State Management (Redis)**: Maintaining a high-speed, volatile record of active challenges and solved sessions.

## Architecture and Resilience

The service follows a strict **Layered Pattern: Handler -> Process -> Task**, ensuring high testability and clear separation of concerns.

### Modular Design
1. **Handler**: Manages HTTP/REST communication, unmarshaling requests, and handling status codes.
2. **Process**: The "brain" of the operation. It orchestrates complex flows, such as verifying a PoW solution before triggering the captcha generation.
3. **Task**: Atomic, reusable components performing specific logic, such as `VerifyPowTask`, `GenerateCaptchaTask`, or `SaveUsedSeedTask`.

### Security & Reliability Features
- **HMAC-Signed Seeds**: Every PoW challenge is signed with a server-side secret, making it impossible for clients to forge their own challenges.
- **Infrastructure Retry Strategy**: Automatically waits for Redis to become available during startup, ensuring stability in containerized environments.
- **Context-Aware Execution**: Full `context.Context` integration for precise timeout control and resource management.
- **Minimal Footprint**: Built using multi-stage Docker builds on Alpine Linux, optimized for security and fast deployment.

## Technical Specification

- **Go**: 1.23+ (utilizing advanced concurrency and type safety).
- **Redis**: Primary store for session states, PoW seeds, and anti-double-spend records.
- **HMAC-SHA256**: Used for secure signing of PoW seeds.
- **Base64 Image Streaming**: Captcha images are generated and streamed as Base64 strings for seamless frontend integration.

## Environment Configuration

The service utilizes a strict configuration validation policy to ensure all security parameters are present at runtime.

| Variable    | Description |
|-------------|-------------|
| APP_ENV     | Runtime environment (local/production) |
| REDIS_URL   | Connection string for the Redis state store |
| HMAC_SECRET | Secret key used for HMAC signing of PoW seeds |

## Data Flow: Protection Sequence

1. **Step 1: Challenge Request**: Client requests a PoW seed. The service returns a signed seed and timestamp.
2. **Step 2: Proof Submission**: Client solves the PoW (finds the nonce) and submits it.
3. **Step 3: Validation & Generation**:
    - Process validates the signature, timestamp, and checks for double-spending.
    - If valid, the service generates a visual captcha and saves the answer to Redis.
4. **Step 4: Final Verification**: User submits the visual answer. The service verifies it, decrements tries on failure, or marks the session as solved on success.

---
Adrian Janczenia