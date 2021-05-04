# Vaccine Notifier

I have used go lang to call the API and GitHub actions to run this script periodically. This code is calling the Cowin API to fetch the details related to covid centers by Pincode. After getting the response it is checking if the vaccine slots are available or not. If slots are available the script will send me the push notification on my mobile.

I have used push notifications as I have one android app to receive the notification. You can simply change the `sendPush` function to send the notification on email as well.

I have used the GitHub action to call this script every 5 mins. Please check .github/workflows/vaccine-notifier.yml to understand the GitHub action.

If you want to use this code then clone it replace the Pincode variable and update the sendPush function as per your need to receive the notification.
