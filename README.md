# KENCAT-Express
# Parcel Delivery Management System – API Overview

## 1. Authentication & User
## Responsible: Aida, Natthida Bunsuea, 6609650350

### POST /api/auth/login
**Description:**  
Authenticate an employee and grant access to the Parcel Delivery Management System.

**Details:** 
This endpoint verifies the employee’s login credentials (such as email/username and password).  
The system validates the credentials against the stored user records in the database. If the authentication is successful, the system returns an access token that allows the employee to access protected system resources.

The authentication process may include password verification using secure hashing methods and validation of the user’s account status.

**Purpose:**  
- Authenticate employees before accessing the system  
- Ensure only authorized personnel can use system functions  
- Provide a secure access token for subsequent API requests  
- Protect the system from unauthorized access


### GET /api/users/{userId}
**Description:**  
Retrieve detailed information of a specific user.

**Details:**  
This endpoint returns the profile information of a user based on the provided `userId`.  
The returned data may include personal information such as name, contact information, and other relevant user details stored in the system.

**Purpose:**  
- Allow employees to view user profile information  
- Support customer service and parcel management


---

## 2. Parcel Management

### POST /api/parcels
**Description:**  
Create a new parcel record in the system.

**Details:**  
This endpoint is used by the **Parcel Clerk** to register a new parcel that has been received from a sender.  
The system stores parcel details such as sender information, receiver information, parcel weight, and delivery type.

**Purpose:**  
- Register parcels in the system  
- Prepare parcels for delivery and tracking


### GET /api/parcels
**Description:**  
Retrieve a list of all parcels in the system.

**Details:**  
This endpoint returns a collection of parcels currently stored in the database.  
It allows employees to view and manage parcels that are pending, in transit, or already delivered.

**Purpose:**  
- Provide parcel overview for employees  
- Support parcel management operations


---

## 3. Parcel Details & Shipping Calculation
## Responsible: Gus, Natchathorn Danusawad, 6609650319

### GET /api/parcels/{parcelId}
**Description:**  
Retrieve detailed information about a specific parcel.

**Details:**  
This endpoint returns complete parcel information including sender, receiver, weight, delivery type, and tracking reference.

**Purpose:**  
- View detailed parcel information  
- Support tracking and delivery management


### GET /api/shipping/calculate
**Description:**  
Calculate the estimated shipping cost.

**Details:**  
This endpoint calculates the shipping fee based on input parameters such as parcel weight, delivery type, and distance.  
The system processes these parameters and returns the estimated cost for the delivery service.

**Purpose:**  
- Allow users to estimate shipping cost before creating a shipment  
- Ensure transparent pricing for customers


---

## 4. Parcel Tracking
## Responsible: Muay, Piyatida Reakdee, 6609650491

### GET /api/trackings/{trackId}
**Description:**  
Retrieve the tracking status of a parcel.

**Details:**  
This endpoint allows users or employees to check the current status of a parcel using the tracking ID.  
The response includes delivery status, timestamps, and other tracking-related information.

**Purpose:**  
- Enable real-time parcel tracking  
- Provide delivery transparency for customers


### PUT /api/trackings/{trackId}/status
*(or PATCH depending on implementation)*

**Description:**  
Update the delivery status of a parcel.

**Details:**  
This endpoint is used by **Messenger staff** to update the delivery status of a parcel during the delivery process.  
Possible statuses may include:

- Pending  
- In Transit  
- Delivered  
- Delivery Failed

Additional information such as delivery confirmation photos may also be stored.

**Purpose:**  
- Maintain accurate delivery status updates  
- Support real-time tracking information


---

## 5. Messenger Operations

### GET /api/messenger/tasks
**Description:**  
Retrieve delivery tasks assigned to a messenger.

**Details:**  
This endpoint returns a list of parcels assigned to a messenger for delivery.  
The messenger can view delivery details and prioritize tasks accordingly.

**Purpose:**  
- Allow messengers to view their assigned delivery tasks  
- Support route planning and delivery operations


### POST /api/vehicle/assign
**Description:**  
Assign a vehicle for delivery.

**Details:**  
This endpoint allows a messenger to select or assign a vehicle to be used for parcel delivery.  
The system records the selected vehicle and associates it with the delivery operation.

**Purpose:**  
- Track vehicle usage during delivery  
- Support transportation management


---

## 6. Reporting System
## Responsible: Gin, Noppawat Jamjumrus, 6609650301

### POST /api/reports
**Description:**  
Submit a report for a parcel-related issue.

**Details:**  
This endpoint allows users to report problems related to a parcel, such as damage, loss, or delivery issues.  
The system records the report along with parcel and user references.

**Purpose:**  
- Provide a mechanism for reporting delivery issues  
- Support issue resolution and customer service


### GET /api/reports
**Description:**  
Retrieve all reported parcel issues.

**Details:**  
This endpoint returns a list of reports submitted by users.  
Employees can review these reports to investigate and resolve delivery problems.

**Purpose:**  
- Allow administrators or staff to monitor reported issues  
- Improve service quality and problem resolution
