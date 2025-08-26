# Project Backend 04 — Go_Bootcamp

**Summary:**  
In this project, you'll learn how to implement and manage authorization and a database in a web application written in Go, using net/http and jackc/pgx.

💡 [Click here](https://new.oprosso.net/p/4cb31ec3f47a4596bc758ea1861fb624) to share your feedback on this project. It's anonymous and will help our team improve the learning process. We recommend completing the survey right after finishing the project.

## Contents

  - [Chapter I](#chapter-i)
    - [Instructions](#instructions)
  - [Chapter II](#chapter-ii)
    - [General Information](#general-information)
      - [Authorization](#authorization)
      - [Identification, Authentication, Authorization](#identification-authentication-authorization)
      - [Authorization Using Login and Password](#authorization-using-login-and-password)
    - [Topics for Study:\*\*](#topics-for-study)
  - [Chapter III](#chapter-iii)
  - [Project: Tic-Tac-Toe](#project-tic-tac-toe)
    - [Task 1. Adding a database](#task-1-adding-a-database)
    - [Task 2. Adding authorization](#task-2-adding-authorization)
    - [Task 3. Adding logic for multiplayer games](#task-3-adding-logic-for-multiplayer-games)


## Chapter I

### Instructions

1. Throughout the course, you’ll experience uncertainty and a lack of information — this is normal. Remember, the repository, Google, your peers, and Rocket.Chat are always available. Communicate. Search. Use common sense. Don’t be afraid to make mistakes.
2. Pay close attention to your sources. Verify. Think critically. Analyze. Compare.
3. Read each task carefully. Then read it again.
4. It’s best to read examples attentively as well. They may contain things not explicitly mentioned in the task.
5. You may encounter inconsistencies — something in the instructions or an example that contradicts what you previously learned. If this happens, try to understand it. If that fails, write down your question and resolve it later. Don’t leave open questions unresolved.
6. If a task seems confusing or impossible — it only seems that way. Try breaking it down. Most parts will likely become clear.
7. You’ll come across a variety of tasks. Those marked with an asterisk (\*) are optional and more challenging. Completing them will provide you with extra experience and knowledge.
8. Don’t try to trick the system or others. In the end, you're only fooling yourself.
9. Have a question? Ask your neighbor on the right. If that doesn’t help, ask the one on the left.
10. Whenever you receive help, always make sure you understand **why**, **how**, and **what for**. Otherwise, the help is pointless.
11. Always push only to the develop branch! The master branch will be ignored. Work in the src directory.
12. Your directory should contain only the files specified in the tasks — no more.

## Chapter II

### General Information

#### Authorization

Authorization mechanisms control access to system resources by legitimate users, assigning them the exact permissions defined by the system administrator.

#### Identification, Authentication, Authorization

- **Identification** — a process by which a system determines a unique trait that unambiguously identifies a subject within the system.
- **Authentication** — a process for verifying authenticity, such as confirming a user's identity by comparing the entered password with the one stored in the system.
- **Authorization** — the process of granting an individual or group the right to perform a specific set of actions.

#### Authorization Using Login and Password

This method requires the user to provide a login and password for successful identification and authentication in the system. The login-password pair is set by the user during registration. Upon successful authorization, the server grants the user permission to perform available requests.

The client sends a request to the server and receives a response with the message “Unauthorized” along with information about the authorization process. Once authorization is successfully completed, each subsequent request automatically includes the Authorization header ([header formation](https://datatracker.ietf.org/doc/html/rfc7617)), which transmits the client’s data for authentication by the server.

![img](misc/images/Auth.png/)

Other [authorization methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#authentication_schemes) also exist.

### Topics for Study:**

- Web applications
- Authorization via login-password pair (basic auth)
- PostgreSQL
- jackc/pgx
- net/http

## Chapter III

## Project: Tic-Tac-Toe

Use the backend project from the previous week (T03) as the basis.

### Task 1. Adding a database

- Set up a connection to a PostgreSQL database using the jackc/pgx library.
- Remove the in-memory storage structure.
- Add struct tags to the fields that need to be persisted in the database.

### Task 2. Adding authorization

- Add support for users, each having a UUID, login, and password.
- Implement user support across all layers of the application.
- Create a SignUpRequest model that includes login and password.
- Create an authentication service that uses UserService with the following methods:
  - A registration method that takes a SignUpRequest and returns an indicator of successful registration.
  - An authentication method that receives login and password in the header encoded as base64(login:password) and returns the user’s UUID.
- Create an authorization handler with the following endpoints:
  - user registration
  - user authentication
- Create a UserAuthenticator structure that protects endpoints from being accessed by unauthorized users:
  - Validate the login and password.
  - If validation succeeds, allow the request.
  - If validation fails, respond with HTTP 401 and do not process the request.
- Apply the UserAuthenticator to your endpoints:
  - Allow unauthenticated access to the registration and authentication endpoints.
  - Require authentication for all other endpoints.

### Task 3. Adding logic for multiplayer games

- Add states for the current game:
  - Waiting for players
  - Player with UUID to move
  - Draw
  - Player with UUID wins
- Add information about the symbols used by each player in the current game.
- Improve the end-of-game detection algorithm using game states.
- Add an endpoint for creating a new game against a player or the computer.
- Add an endpoint for retrieving available ongoing games.
- Add an endpoint for joining a game as a player.
- Improve the endpoint for updating the current game to support both player-vs-player and player-vs-computer modes.
- Add an endpoint for retrieving the current game.
- Add an endpoint for retrieving user information by UUID.