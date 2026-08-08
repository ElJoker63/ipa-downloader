export namespace models {
	
	export class AccountProfile {
	    name: string;
	    email: string;
	    storeFront: string;
	    storeFrontCountry: string;
	    directoryServicesId: string;
	    pod: string;
	    isLoggedIn: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AccountProfile(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.email = source["email"];
	        this.storeFront = source["storeFront"];
	        this.storeFrontCountry = source["storeFrontCountry"];
	        this.directoryServicesId = source["directoryServicesId"];
	        this.pod = source["pod"];
	        this.isLoggedIn = source["isLoggedIn"];
	    }
	}
	export class VersionInfo {
	    externalVersionId: string;
	    displayVersion: string;
	    // Go type: time
	    releaseDate: any;
	    formattedDate: string;
	
	    static createFrom(source: any = {}) {
	        return new VersionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.externalVersionId = source["externalVersionId"];
	        this.displayVersion = source["displayVersion"];
	        this.releaseDate = this.convertValues(source["releaseDate"], null);
	        this.formattedDate = source["formattedDate"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class AppMetadata {
	    id: number;
	    bundleId: string;
	    name: string;
	    developer: string;
	    developerId?: number;
	    version: string;
	    buildNumber?: string;
	    price: number;
	    formattedPrice: string;
	    currency?: string;
	    artworkUrl: string;
	    artworkUrl512?: string;
	    screenshots: string[];
	    ipadScreenshots?: string[];
	    description: string;
	    releaseNotes?: string;
	    releaseDate: string;
	    currentVersionDate?: string;
	    minimumOsVersion: string;
	    fileSizeBytes: number;
	    formattedSize: string;
	    averageUserRating: number;
	    userRatingCount: number;
	    contentAdvisoryRating: string;
	    genres: string[];
	    primaryGenre: string;
	    supportedPlatforms: string[];
	    isFavorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppMetadata(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.bundleId = source["bundleId"];
	        this.name = source["name"];
	        this.developer = source["developer"];
	        this.developerId = source["developerId"];
	        this.version = source["version"];
	        this.buildNumber = source["buildNumber"];
	        this.price = source["price"];
	        this.formattedPrice = source["formattedPrice"];
	        this.currency = source["currency"];
	        this.artworkUrl = source["artworkUrl"];
	        this.artworkUrl512 = source["artworkUrl512"];
	        this.screenshots = source["screenshots"];
	        this.ipadScreenshots = source["ipadScreenshots"];
	        this.description = source["description"];
	        this.releaseNotes = source["releaseNotes"];
	        this.releaseDate = source["releaseDate"];
	        this.currentVersionDate = source["currentVersionDate"];
	        this.minimumOsVersion = source["minimumOsVersion"];
	        this.fileSizeBytes = source["fileSizeBytes"];
	        this.formattedSize = source["formattedSize"];
	        this.averageUserRating = source["averageUserRating"];
	        this.userRatingCount = source["userRatingCount"];
	        this.contentAdvisoryRating = source["contentAdvisoryRating"];
	        this.genres = source["genres"];
	        this.primaryGenre = source["primaryGenre"];
	        this.supportedPlatforms = source["supportedPlatforms"];
	        this.isFavorite = source["isFavorite"];
	    }
	}
	export class AppDetailsOutput {
	    metadata: AppMetadata;
	    versionHistory: VersionInfo[];
	    isFavorite: boolean;
	
	    static createFrom(source: any = {}) {
	        return new AppDetailsOutput(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.metadata = this.convertValues(source["metadata"], AppMetadata);
	        this.versionHistory = this.convertValues(source["versionHistory"], VersionInfo);
	        this.isFavorite = source["isFavorite"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class AppSettings {
	    theme: string;
	    language: string;
	    defaultDownloadFolder: string;
	    maxConcurrentDownloads: number;
	    autoCheckUpdates: boolean;
	    autoAcquireLicense: boolean;
	    rememberCredentials: boolean;
	    keychainPassphrase?: string;
	    searchLimit: number;
	
	    static createFrom(source: any = {}) {
	        return new AppSettings(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.theme = source["theme"];
	        this.language = source["language"];
	        this.defaultDownloadFolder = source["defaultDownloadFolder"];
	        this.maxConcurrentDownloads = source["maxConcurrentDownloads"];
	        this.autoCheckUpdates = source["autoCheckUpdates"];
	        this.autoAcquireLicense = source["autoAcquireLicense"];
	        this.rememberCredentials = source["rememberCredentials"];
	        this.keychainPassphrase = source["keychainPassphrase"];
	        this.searchLimit = source["searchLimit"];
	    }
	}
	export class DeviceInfo {
	    udid: string;
	    name: string;
	    model: string;
	    productType: string;
	    deviceClass: string;
	    iosVersion: string;
	    buildVersion: string;
	    serialNumber: string;
	    wifiAddress?: string;
	    storageTotal: number;
	    storageUsed: number;
	    storageFree: number;
	    batteryLevel: number;
	    batteryCharging: boolean;
	    isPaired: boolean;
	    isConnected: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DeviceInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.udid = source["udid"];
	        this.name = source["name"];
	        this.model = source["model"];
	        this.productType = source["productType"];
	        this.deviceClass = source["deviceClass"];
	        this.iosVersion = source["iosVersion"];
	        this.buildVersion = source["buildVersion"];
	        this.serialNumber = source["serialNumber"];
	        this.wifiAddress = source["wifiAddress"];
	        this.storageTotal = source["storageTotal"];
	        this.storageUsed = source["storageUsed"];
	        this.storageFree = source["storageFree"];
	        this.batteryLevel = source["batteryLevel"];
	        this.batteryCharging = source["batteryCharging"];
	        this.isPaired = source["isPaired"];
	        this.isConnected = source["isConnected"];
	    }
	}
	export class DownloadTask {
	    id: string;
	    appId: number;
	    bundleId: string;
	    appName: string;
	    developer: string;
	    version: string;
	    artworkUrl: string;
	    destinationPath: string;
	    status: string;
	    totalBytes: number;
	    downloadedBytes: number;
	    progress: number;
	    speedBytesPerSec: number;
	    formattedSpeed: string;
	    etaSeconds: number;
	    formattedETA: string;
	    externalVersionId?: string;
	    platform: string;
	    error?: string;
	    // Go type: time
	    createdAt: any;
	    // Go type: time
	    updatedAt: any;
	    // Go type: time
	    completedAt?: any;
	
	    static createFrom(source: any = {}) {
	        return new DownloadTask(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.appId = source["appId"];
	        this.bundleId = source["bundleId"];
	        this.appName = source["appName"];
	        this.developer = source["developer"];
	        this.version = source["version"];
	        this.artworkUrl = source["artworkUrl"];
	        this.destinationPath = source["destinationPath"];
	        this.status = source["status"];
	        this.totalBytes = source["totalBytes"];
	        this.downloadedBytes = source["downloadedBytes"];
	        this.progress = source["progress"];
	        this.speedBytesPerSec = source["speedBytesPerSec"];
	        this.formattedSpeed = source["formattedSpeed"];
	        this.etaSeconds = source["etaSeconds"];
	        this.formattedETA = source["formattedETA"];
	        this.externalVersionId = source["externalVersionId"];
	        this.platform = source["platform"];
	        this.error = source["error"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	        this.updatedAt = this.convertValues(source["updatedAt"], null);
	        this.completedAt = this.convertValues(source["completedAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class FavoriteApp {
	    appId: number;
	    bundleId: string;
	    name: string;
	    developer: string;
	    version: string;
	    price: number;
	    formattedPrice: string;
	    artworkUrl: string;
	    primaryGenre: string;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new FavoriteApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.appId = source["appId"];
	        this.bundleId = source["bundleId"];
	        this.name = source["name"];
	        this.developer = source["developer"];
	        this.version = source["version"];
	        this.price = source["price"];
	        this.formattedPrice = source["formattedPrice"];
	        this.artworkUrl = source["artworkUrl"];
	        this.primaryGenre = source["primaryGenre"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class IPAInfo {
	    bundleId: string;
	    bundleName: string;
	    version: string;
	    shortVersion: string;
	    minimumOs: string;
	    architectures: string[];
	    fileSizeBytes: number;
	    isValid: boolean;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new IPAInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.bundleId = source["bundleId"];
	        this.bundleName = source["bundleName"];
	        this.version = source["version"];
	        this.shortVersion = source["shortVersion"];
	        this.minimumOs = source["minimumOs"];
	        this.architectures = source["architectures"];
	        this.fileSizeBytes = source["fileSizeBytes"];
	        this.isValid = source["isValid"];
	        this.error = source["error"];
	    }
	}
	export class InstalledApp {
	    name: string;
	    bundleId: string;
	    version: string;
	    shortVersion: string;
	    size: number;
	    dynamicSize?: number;
	    vendor?: string;
	    appType: string;
	    minimumOs?: string;
	    signerIdentity?: string;
	
	    static createFrom(source: any = {}) {
	        return new InstalledApp(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.bundleId = source["bundleId"];
	        this.version = source["version"];
	        this.shortVersion = source["shortVersion"];
	        this.size = source["size"];
	        this.dynamicSize = source["dynamicSize"];
	        this.vendor = source["vendor"];
	        this.appType = source["appType"];
	        this.minimumOs = source["minimumOs"];
	        this.signerIdentity = source["signerIdentity"];
	    }
	}
	export class LogEntry {
	    id: number;
	    // Go type: time
	    timestamp: any;
	    level: string;
	    message: string;
	    context?: string;
	
	    static createFrom(source: any = {}) {
	        return new LogEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.timestamp = this.convertValues(source["timestamp"], null);
	        this.level = source["level"];
	        this.message = source["message"];
	        this.context = source["context"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class SearchHistoryItem {
	    id: number;
	    term: string;
	    platform: string;
	    count: number;
	    // Go type: time
	    createdAt: any;
	
	    static createFrom(source: any = {}) {
	        return new SearchHistoryItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.term = source["term"];
	        this.platform = source["platform"];
	        this.count = source["count"];
	        this.createdAt = this.convertValues(source["createdAt"], null);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}

}

