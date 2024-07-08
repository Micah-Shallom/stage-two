# Stage-two

## **HNG internship: stage two backend task**

### Step-by-Step Task Explanation

1. **Setup & Database Connection**

   - Choose a backend framework (e.g., Express.js for Node.js, Flask for Python, etc.).
   - Connect the application to a PostgreSQL database.
   - (Optional) Integrate an ORM (e.g., Sequelize for Node.js, SQLAlchemy for Python).
2. **User Model Creation**

   - Define a User model with the following properties:
     ```json
     {
       "userId": "string", // Unique
       "firstName": "string", // Required, not null
       "lastName": "string", // Required, not null
       "email": "string", // Unique, required, not null
       "password": "string", // Required, not null
       "phone": "string"
     }
     ```
   - Ensure `userId` and `email` are unique.
   - Add validation for all fields.
3. **Validation Handling**

   - Implement validation logic.
   - If validation fails, return a 422 status code with:
     ```json
     {
       "errors": [
         {
           "field": "string",
           "message": "string"
         }
       ]
     }
     ```
4. **User Authentication Implementation**

   - **User Registration**

     - Create a registration endpoint (`/auth/register`).
     - Hash the user’s password before storing it in the database.
     - On successful registration, return a 201 status code with:
       ```json
       {
         "status": "success",
         "message": "Registration successful",
         "data": {
           "accessToken": "eyJh...",
           "user": {
             "userId": "string",
             "firstName": "string",
             "lastName": "string",
             "email": "string",
             "phone": "string"
           }
         }
       }
       ```
     - On failure, return a 400 status code with:
       ```json
       {
         "status": "Bad request",
         "message": "Registration unsuccessful",
         "statusCode": 400
       }
       ```
   - **User Login**

     - Create a login endpoint (`/auth/login`).
     - Validate user credentials and generate a JWT token on successful login.
     - Return a 200 status code with:
       ```json
       {
         "status": "success",
         "message": "Login successful",
         "data": {
           "accessToken": "eyJh...",
           "user": {
             "userId": "string",
             "firstName": "string",
             "lastName": "string",
             "email": "string",
             "phone": "string"
           }
         }
       }
       ```
     - On failure, return a 401 status code with:
       ```json
       {
         "status": "Bad request",
         "message": "Authentication failed",
         "statusCode": 401
       }
       ```
5. **Organisation Management**

   - Define an Organisation model with the following properties:
     ```json
     {
       "orgId": "string", // Unique
       "name": "string", // Required, not null
       "description": "string"
     }
     ```
6. **Endpoints**

   - **[POST] /auth/register**

     - Registers a user and creates a default organisation with the user’s first name appended with "Organisation".
   - **[POST] /auth/login**

     - Logs in a user and returns a JWT token.
   - **[GET] /api/users/:id**

     - Retrieves a user’s record (PROTECTED endpoint).
     - Return 200 status code with user data.
   - **[GET] /api/organisations**

     - Retrieves all organisations the logged-in user belongs to or created (PROTECTED endpoint).
     - Return 200 status code with organisation data.
   - **[GET] /api/organisations/:orgId**

     - Retrieves a single organisation record (PROTECTED endpoint).
     - Return 200 status code with organisation data.
   - **[POST] /api/organisations**

     - Allows a user to create a new organisation (PROTECTED endpoint).
     - Return 201 status code with organisation data on success, 400 on failure.
   - **[POST] /api/organisations/:orgId/users**

     - Adds a user to a specific organisation (PROTECTED endpoint).
     - Return 200 status code on success.
7. **Unit Testing**

   - **Token Generation**

     - Ensure token expiration and correct user details in the token.
   - **Organisation Access**

     - Ensure users can’t see data from organisations they don’t have access to.
8. **End-to-End Testing for Register Endpoint**

   - **Test File Structure**

     - Name the test file `auth.spec.ext` inside a `tests` folder.
   - **Test Scenarios**

     - Successful user registration with default organisation.
     - Successful user login with valid credentials.
     - Registration failure for missing required fields.
     - Registration failure for duplicate email or userId.

   # **Unit Testing**

   Write appropriate unit tests to cover

   Token generation - Ensure token expires at the correct time and correct user details is found in token.
   Organisation - Ensure users can’t see data from organisations they don’t have access to.
   ## **End-to-End Test Requirements for the Register Endpoint**
   The goal is to ensure the POST /auth/register endpoint works correctly by performing end-to-end tests. The tests should cover successful user registration, validation errors, and database constraints.
   Directory Structure:

   The test file should be named auth.spec.ext (ext is the file extension of your chosen language) inside a folder named tests . For example tests/auth.spec.ts assuming I’m using Typescript
   
   ### **Test Scenarios:**
   -    It Should Register User Successfully with Default Organisation:Ensure a user is registered successfully when no organisation details are provided.
   -    Verify the default organisation name is correctly generated (e.g., "John's Organisation" for a user with the first name "John").
   -    Check that the response contains the expected user details and access token.
   -    It Should Log the user in successfully:Ensure a user is logged in successfully when a valid credential is provided and fails otherwise.
   -    Check that the response contains the expected user details and access token.
   -    It Should Fail If Required Fields Are Missing:Test cases for each required field (firstName, lastName, email, password) missing.
   -    Verify the response contains a status code of 422 and appropriate error messages.
   -    It Should Fail if there’s Duplicate Email or UserID:Attempt to register two users with the same email.
   -    Verify the response contains a status code of 422 and appropriate error messages.