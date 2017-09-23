// This is modified from the ios implementation in github.com/SpiderOak/keyring

package keyring

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework Foundation -framework Security
#import <Foundation/Foundation.h>
#import <Security/Security.h>
int Set(char *service, char *username, char *password) {
	NSMutableDictionary *keychainItem = [NSMutableDictionary dictionary];
	keychainItem[(__bridge id)kSecClass] = (__bridge id)kSecClassGenericPassword;
	keychainItem[(__bridge id)kSecAttrAccessible] = (__bridge id)kSecAttrAccessibleWhenUnlocked;
	keychainItem[(__bridge id)kSecAttrAccount] = [NSString stringWithUTF8String:username];
	keychainItem[(__bridge id)kSecAttrService] = [NSString stringWithUTF8String:service];
	if(SecItemCopyMatching((__bridge CFDictionaryRef)keychainItem, NULL) == noErr) {
		NSMutableDictionary *attributesToUpdate = [NSMutableDictionary dictionary];
		attributesToUpdate[(__bridge id)kSecValueData] = [@(password) dataUsingEncoding:NSUTF8StringEncoding];
        OSStatus sts = SecItemUpdate((__bridge CFDictionaryRef)keychainItem, (__bridge CFDictionaryRef)attributesToUpdate);
		return (int)sts;
	} else {
		keychainItem[(__bridge id)kSecValueData] = [@(password) dataUsingEncoding:NSUTF8StringEncoding];
		OSStatus sts = SecItemAdd((__bridge CFDictionaryRef)keychainItem, NULL);
		return (int)sts;
	}
}
const char *Get(char *service, char *username) {
	NSMutableDictionary *keychainItem = [NSMutableDictionary dictionary];
	keychainItem[(__bridge id)kSecClass] = (__bridge id)kSecClassGenericPassword;
	keychainItem[(__bridge id)kSecAttrAccessible] = (__bridge id)kSecAttrAccessibleWhenUnlocked;
	keychainItem[(__bridge id)kSecAttrAccount] = [NSString stringWithUTF8String:username];
	keychainItem[(__bridge id)kSecAttrService] = [NSString stringWithUTF8String:service];
	keychainItem[(__bridge id)kSecReturnData] = (__bridge id)kCFBooleanTrue;
	keychainItem[(__bridge id)kSecReturnAttributes] = (__bridge id)kCFBooleanTrue;
	CFDictionaryRef result = nil;
	OSStatus sts = SecItemCopyMatching((__bridge CFDictionaryRef)keychainItem, (CFTypeRef *)&result);
	if(sts == noErr) {
		NSDictionary *resultDict = (NSDictionary *)result;
		NSData *pswd = resultDict[(__bridge id)kSecValueData];
		return [pswd bytes];
	} else {
		NSLog(@"Keychain Get: Error Code: %d", (int)sts);
	}
	return nil;
}
*/
import "C"

import (
	"fmt"
)

type osxProvider struct {
}

func init() {
	defaultProvider = osxProvider{}
}

func (p osxProvider) Get(Service, Username string) (string, error) {
	res := C.GoString(C.Get(C.CString(Service), C.CString(Username)))
	if len(res) == 0 {
		return "", ErrNotFound
	}

	return res, nil
}

const (
	errSecSuccess               = 0      /* No error. */
	errSecUnimplemented         = -4     /* Function or operation not implemented. */
	errSecParam                 = -50    /* One or more parameters passed to a function where not valid. */
	errSecAllocate              = -108   /* Failed to allocate memory. */
	errSecNotAvailable          = -25291 /* No keychain is available. You may need to restart your computer. */
	errSecDuplicateItem         = -25299 /* The specified item already exists in the keychain. */
	errSecItemNotFound          = -25300 /* The specified item could not be found in the keychain. */
	errSecInteractionNotAllowed = -25308 /* User interaction is not allowed. */
	errSecDecode                = -26275 /* Unable to decode the provided data. */
	errSecAuthFailed            = -25293 /* The user name or passphrase you entered is not correct. */
)

func (p osxProvider) Set(Service, Username, Password string) error {
	ret := C.Set(C.CString(Service), C.CString(Username), C.CString(Password))
	switch ret {
	case errSecSuccess:
		return nil
	case errSecUnimplemented:
		return fmt.Errorf("Function or operation not implemented")
	case errSecParam:
		return fmt.Errorf("One or more parameters passed to a function where not valid")
	case errSecAllocate:
		return fmt.Errorf("Failed to allocate memory")
	case errSecNotAvailable:
		return fmt.Errorf("No keychain is available. You may need to restart your computer")
	case errSecDuplicateItem:
		return fmt.Errorf("The specified item already exists in the keychain")
	case errSecItemNotFound:
		return fmt.Errorf("The specified item could not be found in the keychain")
	case errSecInteractionNotAllowed:
		return fmt.Errorf("User interaction is not allowed")
	case errSecDecode:
		return fmt.Errorf("Unable to decode the provided data")
	case errSecAuthFailed:
		return fmt.Errorf("The user name or passphrase you entered is not correct")
	default:
		return fmt.Errorf("Unknown return code from iOS: %d", ret)
	}
	return nil
}


