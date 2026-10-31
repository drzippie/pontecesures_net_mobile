# Add project specific ProGuard rules here.
# You can control the set of applied configuration files using the
# proguardFiles setting in build.gradle.

# Keep native methods
-keepclasseswithmembernames class * {
    native <methods>;
}

# Keep Wails bridge classes
-keep class com.wails.app.WailsBridge { *; }
-keep class com.wails.app.WailsJSBridge { *; }

# Methods exposed to the page through addJavascriptInterface ("wails" object).
-keepattributes JavascriptInterface
-keepclassmembers class * {
    @android.webkit.JavascriptInterface <methods>;
}

# Go calls back into Java through JNI by method name on the WailsBridge
# instance (GetObjectClass + GetMethodID: executeJavaScript, share, openURL,
# secureSet, ...). The class is kept whole above; never obfuscate it.
-keepnames class com.wails.app.** { *; }

# Compile-time-only annotations referenced by Tink (androidx.security-crypto).
-dontwarn javax.annotation.Nullable
-dontwarn javax.annotation.concurrent.GuardedBy
