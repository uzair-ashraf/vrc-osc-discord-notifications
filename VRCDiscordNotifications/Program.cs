using System;
using System.Linq;
using System.Threading;
using System.Threading.Tasks;
using Windows.UI.Notifications.Management;
using Windows.Foundation.Metadata;
using System.Collections.Generic;
using Windows.UI.Notifications;
using System.Text.Json;

namespace VRCDiscordNotifications {
    class Program {
        private static HashSet<uint> seenIDs = new HashSet<uint>();
        static void Main(string[] args) {
            if (!ApiInformation.IsTypePresent("Windows.UI.Notifications.Management.UserNotificationListener")) {
                Console.WriteLine(JsonSerializer.Serialize(new Dictionary<string, bool>(){ ["WasPermissionGranted"] = false }));
                Console.ReadLine();
                return;
            }
            Task.Run(async () => {  
                bool wasPermissionGranted = await IsNotificationPermissionsGranted();
                if (!wasPermissionGranted) {
                    Console.WriteLine(JsonSerializer.Serialize(new Dictionary<string, bool>(){ ["WasPermissionGranted"] = false }));
                    Console.ReadLine();
                    return;
                }
                Console.WriteLine(JsonSerializer.Serialize(new Dictionary<string, bool>() { ["WasPermissionGranted"] = true }));
                    while (true) {
                    UserNotificationListener listener = UserNotificationListener.Current;
                    IReadOnlyList<UserNotification> notifs = await listener.GetNotificationsAsync(NotificationKinds.Toast);
                    foreach (UserNotification n in notifs) {
                        NotificationBinding toastBinding = n.Notification.Visual.GetBinding(KnownNotificationBindings.ToastGeneric);
                        if (toastBinding == null) continue;
                        string appName = n.AppInfo.DisplayInfo.DisplayName;
                        if (appName != "Discord") continue; 
                        IReadOnlyList<AdaptiveNotificationText> textElements = toastBinding.GetTextElements();
                        string username = textElements.FirstOrDefault()?.Text;
                        if (username == null || username == "") {
                            continue;
                        }
                        if (seenIDs.Contains(n.Id)) {
                            continue;
                        }
                        Console.WriteLine(JsonSerializer.Serialize(new Dictionary<string, string>() { ["Username"] = username }));
                        seenIDs.Add(n.Id);
                        Thread.Sleep(1000);
                    }
                }
            }).GetAwaiter().GetResult();
        }
        static async Task<bool> IsNotificationPermissionsGranted() {
            try {
                UserNotificationListener listener = UserNotificationListener.Current;
                UserNotificationListenerAccessStatus accessStatus = await listener.RequestAccessAsync();
                switch (accessStatus) {
                    case UserNotificationListenerAccessStatus.Allowed:
                        return true;
                    case UserNotificationListenerAccessStatus.Denied:
                        return false;
                    case UserNotificationListenerAccessStatus.Unspecified:
                        return false;
                }
                return false;
            } catch (Exception error) {
                Console.WriteLine(error);
                return false;
            }
        }
    }
}
