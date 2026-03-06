# WhatsApp API Documentation for Foonte

Based on the provided links, here is the documentation for the Foonte WhatsApp API:

## 1. Send Message API

The Send Message API allows you to send various types of messages to WhatsApp contacts.

### Endpoint
```
POST https://api.fonnte.com/send
```

### Headers
```
Authorization: {{token}}
```

### Request Body Parameters

#### Required Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| target | String | Yes | Target number, comma separated, variable supported. Can be WhatsApp number, Group ID, or Rotator ID. Format: '081xxxx' or '081xxxx,082xxxx,083xxxx,123xxxx@g.us' |

#### Optional Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| message | String | No | Text message, support variable. Cannot exceed 60,000 characters |
| url | String | No | Public URL of attachment (image, file, audio, video). Note: Only available on super/advanced/ultra packages |
| filename | String | No | Custom filename for non-image/video files |
| schedule | Integer | No | Unix timestamp to send messages on a schedule |
| delay | String | No | Min 0, add - to random delay (for example: 1-10 will produce a random delay between 1 and 10 seconds). Only works on multiple targets |
| countryCode | String | No | Replace first zero with country code, default 62, set 0 to disable replacement |
| location | String | No | Location coordinates in format: latitude,longitude (e.g., '-7.983908, 112.621391') |
| typing | Boolean | No | Typing indicator, default: false |
| choices | String | No | Polling choices, minimum 2, maximum 12, separated by commas (e.g., 'satu,dua,tiga') |
| select | String | No | Polling selection limit: 'single' or 'multiple', default: 'single' |
| pollname | String | No | Name of the poll to identify the poll |
| file | Binary | No | Upload file directly from local/form. Note: Only available on super/advanced/ultra plan |
| connectOnly | Boolean | No | Use the API only on connected devices, default: true |
| followup | Integer | No | Add seconds before sending the message |
| data | String | No | Combine all requests into one request |
| sequence | Boolean | No | Requests must be executed sequentially, default: false |
| preview | Boolean | No | Links in messages should have a preview, default: true |
| inboxid | Integer | No | ID of the inbox message you want to reply |
| duration | Integer | No | Duration for custom typing duration, default: 1 second |

### Example Request
```json
{
  "target": "6281234567890",
  "message": "Hello, this is a test message!",
  "url": "https://md.fonnte.com/images/wa-logo.png",
  "typing": true,
  "delay": "2"
}
```

### Response
Successful Response:
```json
{
  "detail": "success! message in queue",
  "id": [
    "80367170"
  ],
  "process": "pending",
  "requestid": 2937124,
  "status": true,
  "target": [
    "6282227097005"
  ]
}
```

Failed Response Examples:
```json
// Invalid token
{
  "Status": false,
  "reason": "token invalid",
  "requestid": 2937124
}

// Insufficient quota
{
  "reason": "insufficient quota",
  "status": false,
  "requestid": 2937124
}

// Target invalid
{
  "reason": "target invalid",
  "status": false,
  "requestid": 2937124
}
```

### Important Notes
- TOKEN must be filled with your own token
- TOKENS can be multiple, separated by commas, for example: xxxxxx,yyyyyy
- Multiple TOKENS will work as a rotator, each target will be sent by a random device based on the registered token
- Multiple TOKENS can only be used with the same account
- Authorization does not require a Bearer. You can pass the token directly
- Sending multiple numbers can be done by adding a comma
- Sending variables is supported with the API
- Filename only works on file and audio types
- Consider file limitation rules when sending attachments

## 2. Typing Indicator API

The Typing API enables you to simulate typing and create typing indicators on WhatsApp. This is beneficial for adding presence that you are thinking/answering the message or just simulating typing, especially for long-awaited answers like waiting for AI to respond or when your webhook needs to process something that takes longer than a second.

### Endpoint
```
POST https://api.fonnte.com/typing
```

### Headers
```
Authorization: {{token}}
```

### Request Body Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| target | String | Yes | Number to simulate typing with |
| countryCode | String | No | Country code of the target, default 62 |
| duration | Integer | Yes | Typing duration in seconds |
| stop | Boolean | No | Should the typing stop, default: false |

### Example Request
```json
{
  "target": "6281234567890",
  "countryCode": "62",
  "duration": 10,
  "stop": false
}
```

### Response
The API returns a JSON object indicating the success or failure of the typing simulation request.

## 3. Reschedule Message API

This API is used to reschedule messages that have been wrongly scheduled.

### Endpoint
```
POST https://api.fonnte.com/reset-message
```

### Headers
```
Authorization: {{token}}
```

### Request Body Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | Integer | Yes | Message ID |
| delay | String | No | Leave empty if you don't want to change delay. Example: "5-10" |
| schedule | Integer | No | Leave empty if you don't want to change schedule. Unix timestamp format. Example: 1758606641 |
| byschedule | Boolean | No | Affect all messages that have the same schedule. Default: false |

### Example Request
```json
{
  "id": 1,
  "delay": "5-10",
  "schedule": 1758606641,
  "byschedule": false
}
```

### Response
The API returns a JSON object indicating the success or failure of the reschedule request.

### Notes
- The `delay` parameter is in the format "min-max" (e.g., "5-10") representing seconds.
- The `schedule` parameter uses Unix timestamp format for scheduling time.
- When `byschedule` is set to true, all messages with the same schedule will be affected.

## 4. Delete Message API

This API is used to delete a message via API. You can delete the requested message you don't want to proceed anymore or for other reasons, EXCEPT messages with status "processing". Processing messages cannot be deleted. If you want to delete a processing message, you have to disconnect your device first, which will change the processing message to pending, allowing you to delete it.

### Endpoint
```
POST https://api.fonnte.com/delete-message
```

### Headers
```
Authorization: {{token}}
```

### Request Body Parameters
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| id | String | Yes | The ID of the message to be deleted |

### Example Request
```json
{
  "id": "1"
}
```

### Response

#### Success Response
```json
{
  "detail": "message 1 successfully deleted",
  "status": true
}
```

#### Error Responses

##### Invalid Token
```json
{
  "reason": "invalid token",
  "status": false
}
```

##### Invalid Message ID
```json
{
  "reason": "invalid message id",
  "status": false
}
```

##### Cannot Delete Processing Message
```json
{
  "reason": "cannot delete message with status processing",
  "status": false
}
```

### Important Notes
- Messages with status "processing" cannot be deleted directly
- To delete a processing message, disconnect your device first to change its status to "pending"
- Only messages belonging to your device can be deleted
- Requires valid device token for authentication

## 5. Webhook: Reply Message

Webhook makes responding to incoming messages using custom data possible.

**Special note**: Any autoreply feature won't work if you are using webhook.

### Webhook Payload
```json
{
  "device": "device_number",
  "sender": "6281234567890",
  "message": "Hello, this is a reply!",
  "text": "button text",
  "member": "group_member_info",
  "name": "sender_name",
  "location": "latitude,longitude",
  "pollname": "poll_name",
  "choices": "selected_poll_choices",
  "timestamp": 1678886400,
  "inboxid": "inbox_id",
  "url": "attachment_url",
  "filename": "attachment_filename",
  "extension": "file_extension"
}
```

### Fields Explanation
- `device`: Your device number (not connected device)
- `sender`: Sender's WhatsApp number
- `message`: The message content
- `text`: Button text message
- `member`: Group member who sent the message (for group messages)
- `name`: The sender's name
- `location`: Sender's location (latitude,longitude)
- `pollname`: Poll's name
- `choices`: Selected poll choices
- `timestamp`: The timestamp of message received
- `inboxid`: The inbox ID for replying message on the API
- `url`: The file URL if sender sent an attachment (only available on devices with all feature package)
- `filename`: The filename of the attachment (only available on devices with all feature package)
- `extension`: The extension of the attachment (only available on devices with all feature package)

**Note**: Replying using attachment is only available on devices with a package of super, advanced or ultra.

## 6. Webhook: Get Attachment

You can download the attachment sent to your device using webhook. This function will only work on devices with the all-feature package.

**Important Note:** Any autoreply feature won't work if you are using webhook.

### Webhook Payload
```json
{
  "device": "device_identifier",
  "sender": "6281234567890",
  "message": "Attachment message if present",
  "text": "button_text",
  "member": "group_member_info",
  "name": "sender_name",
  "location": "location_data",
  "url": "direct_download_url",
  "filename": "attachment_filename",
  "extension": "file_extension"
}
```

### Fields Explanation

#### Standard Fields (Available for all packages):
- `device`: The device identifier
- `sender`: The sender's phone number or identifier
- `message`: The message content (includes attachment message if present)
- `text`: Button text (if applicable)
- `member`: Group member who sent the message (in group chats)
- `name`: Sender's name
- `location`: Location data (if shared)

#### Special Fields (Only for "all feature" package):
- `url`: Direct URL to download the attachment
- `filename`: Name of the attachment file
- `extension`: File extension of the attachment

### File Limitations

#### Limitations by Type:
- **Image Files**: Supported formats: "png", "jpg", "jpeg", "webp"
- **Video Files**: Supported format: "mp4"
- **Document Files**: Supported formats: "pdf", "doc", "docx", "xls", "xlsx", "csv", "txt"
- **Audio Files**: Supported format: "mp3"
- **Size Limit**: Maximum file size: 10MB

#### System Limitations:
- If you try sending attachments outside these file limitation rules, your API will return false.
- For every attachment outside these file limitation rules sent to your device while using Fonnte, it will not be sent to your webhook or forwarded to the attachment target.
- Any files that don't match the supported formats or exceed the 10MB size limit will not be processed successfully through the Fonnte system.

### Important Notes

1. **Message Handling**: If your attachment has a message with it, you can find it in the `message` field.

2. **File Saving**: The attachment will be downloaded and saved in the same path as your webhook URL. To save elsewhere, specify the path where the attachment should be saved.

3. **Package Requirement**: These special attachment fields (`url`, `filename`, `extension`) are only available for devices with the "all feature" package.

4. **Autoreply Conflict**: Using webhook will disable any autoreply features on your device.

## 7. Webhook: Get Submission

When a client submits their submission, you can receive the data and save it to your own system or process it as needed. The webhook sends a JSON payload containing the submission details.

### Webhook Payload
```json
{
  "submission": "submission_name",
  "sender": "6281234567890",
  "name": "sender_name",
  "data": [
    {
      "question": "answer"
    }
  ]
}
```

### Fields Explanation
- `submission`: The submission name
- `sender`: Sender's WhatsApp number
- `name`: The sender's name
- `data`: The submission data (list of submitted information)

### Usage
You can copy and use this webhook implementation however you like on your system to handle incoming submission data from clients.

## 8. Webhook: Update Message Status

Webhook update message status serves as a replacement for the API message status to enable real-time status updates without requiring an API call to update the status. The message status includes `id` and `stateid` to update the message status and message state.

### Webhook Payload
```json
{
  "device": "device_number",
  "id": "message_id",
  "stateid": "state_id",
  "status": "message_status",
  "state": "message_state"
}
```

### Fields Explanation
- `device`: Your device number (not connected device)
- `id`: The message id
- `stateid`: The message stateid
- `status`: The status of the message
- `state`: The state of the message

### Usage
You need to send from the API to be able to save the report. The webhook can be used to update message status and state in your database in real-time.

## Authentication

All API requests require an Authorization header with your Foonte token:

```
Authorization: {{your_token_here}}
```

## Error Handling

API responses include error information when requests fail:

```json
{
  "success": false,
  "message": "Error description",
  "code": 400
}
```

Common error codes:
- 400: Bad Request (invalid parameters)
- 401: Unauthorized (invalid token)
- 403: Forbidden (insufficient permissions)
- 429: Too Many Requests (rate limiting)
- 500: Internal Server Error

## Rate Limits

Foonte implements rate limiting to prevent abuse. The exact limits depend on your account type and subscription plan. Exceeding rate limits will result in 429 responses.

## Webhook Configuration

To receive webhooks, you need to configure your webhook URL in the Foonte dashboard. The webhook endpoint should be accessible from the internet and respond with a 200 status code within a few seconds to acknowledge receipt.