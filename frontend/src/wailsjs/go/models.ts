export namespace app {
	
	export class ArchiveMobileVideosReq {
	    dirPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ArchiveMobileVideosReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirPath = source["dirPath"];
	    }
	}
	export class ArchiveMobileVideosRes {
	    success: boolean;
	    message: string;
	    movedCount: number;
	    failedFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new ArchiveMobileVideosRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.movedCount = source["movedCount"];
	        this.failedFiles = source["failedFiles"];
	    }
	}
	export class ExtractM3u8LinksReq {
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractM3u8LinksReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	    }
	}
	export class M3u8Info {
	    url: string;
	    title: string;
	
	    static createFrom(source: any = {}) {
	        return new M3u8Info(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.url = source["url"];
	        this.title = source["title"];
	    }
	}
	export class ExtractM3u8LinksRes {
	    success: boolean;
	    message: string;
	    links: M3u8Info[];
	
	    static createFrom(source: any = {}) {
	        return new ExtractM3u8LinksRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.links = this.convertValues(source["links"], M3u8Info);
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
	
	export class MoveFilesReq {
	    sourceFiles: string[];
	    targetDir: string;
	
	    static createFrom(source: any = {}) {
	        return new MoveFilesReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourceFiles = source["sourceFiles"];
	        this.targetDir = source["targetDir"];
	    }
	}
	export class MoveFilesRes {
	    success: boolean;
	    message: string;
	    movedCount: number;
	    failedFiles: string[];
	
	    static createFrom(source: any = {}) {
	        return new MoveFilesRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.movedCount = source["movedCount"];
	        this.failedFiles = source["failedFiles"];
	    }
	}
	export class MovieDownloadReq {
	    mType: string;
	    mid: string;
	    filterStr: string;
	
	    static createFrom(source: any = {}) {
	        return new MovieDownloadReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.mType = source["mType"];
	        this.mid = source["mid"];
	        this.filterStr = source["filterStr"];
	    }
	}
	export class MovieMergeReq {
	    dir: string;
	
	    static createFrom(source: any = {}) {
	        return new MovieMergeReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dir = source["dir"];
	    }
	}
	export class WindowState {
	    X: number;
	    Y: number;
	    Width: number;
	    Height: number;
	
	    static createFrom(source: any = {}) {
	        return new WindowState(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.X = source["X"];
	        this.Y = source["Y"];
	        this.Width = source["Width"];
	        this.Height = source["Height"];
	    }
	}

}

export namespace myVideo {
	
	export class BatchExtractAudioReq {
	    dirPath: string;
	    format: string;
	    threadCount: number;
	
	    static createFrom(source: any = {}) {
	        return new BatchExtractAudioReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirPath = source["dirPath"];
	        this.format = source["format"];
	        this.threadCount = source["threadCount"];
	    }
	}
	export class BatchExtractAudioRes {
	    success: boolean;
	    message: string;
	    totalCount: number;
	    successCount: number;
	    failedCount: number;
	    failedFiles: string[];
	    outputDir: string;
	    totalCost: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchExtractAudioRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.totalCount = source["totalCount"];
	        this.successCount = source["successCount"];
	        this.failedCount = source["failedCount"];
	        this.failedFiles = source["failedFiles"];
	        this.outputDir = source["outputDir"];
	        this.totalCost = source["totalCost"];
	    }
	}
	export class BatchRemoveVideoIntroReq {
	    filePaths: string[];
	    introDuration: number;
	
	    static createFrom(source: any = {}) {
	        return new BatchRemoveVideoIntroReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePaths = source["filePaths"];
	        this.introDuration = source["introDuration"];
	    }
	}
	export class BatchRemoveVideoIntroRes {
	    success: boolean;
	    message: string;
	    processed: number;
	    failed: number;
	    failedFiles: string[];
	    totalCost: string;
	
	    static createFrom(source: any = {}) {
	        return new BatchRemoveVideoIntroRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.processed = source["processed"];
	        this.failed = source["failed"];
	        this.failedFiles = source["failedFiles"];
	        this.totalCost = source["totalCost"];
	    }
	}
	export class ClipAndReplaceReq {
	    filePath: string;
	    startTime: string;
	    endTime: string;
	
	    static createFrom(source: any = {}) {
	        return new ClipAndReplaceReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	    }
	}
	export class ClipAndReplaceRes {
	    success: boolean;
	    message: string;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new ClipAndReplaceRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.cost = source["cost"];
	    }
	}
	export class ClipVideoReq {
	    filePath: string;
	    startTime: string;
	    endTime: string;
	
	    static createFrom(source: any = {}) {
	        return new ClipVideoReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.startTime = source["startTime"];
	        this.endTime = source["endTime"];
	    }
	}
	export class ClipVideoRes {
	    success: boolean;
	    outputPath: string;
	    message: string;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new ClipVideoRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.outputPath = source["outputPath"];
	        this.message = source["message"];
	        this.cost = source["cost"];
	    }
	}
	export class ExtractAudioReq {
	    filePath: string;
	    format: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractAudioReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.format = source["format"];
	    }
	}
	export class ExtractAudioRes {
	    success: boolean;
	    outputPath: string;
	    message: string;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractAudioRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.outputPath = source["outputPath"];
	        this.message = source["message"];
	        this.cost = source["cost"];
	    }
	}
	export class ExtractFramesReq {
	    filePath: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new ExtractFramesReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.count = source["count"];
	    }
	}
	export class ExtractFramesRes {
	    success: boolean;
	    message: string;
	    outputDir: string;
	    framePaths: string[];
	    frameCount: number;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractFramesRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.outputDir = source["outputDir"];
	        this.framePaths = source["framePaths"];
	        this.frameCount = source["frameCount"];
	        this.cost = source["cost"];
	    }
	}
	export class ExtractVideoThumbnailReq {
	    filePath: string;
	    timestamp: number;
	
	    static createFrom(source: any = {}) {
	        return new ExtractVideoThumbnailReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.timestamp = source["timestamp"];
	    }
	}
	export class ExtractVideoThumbnailRes {
	    success: boolean;
	    message: string;
	    thumbnail: string;
	    mimeType: string;
	
	    static createFrom(source: any = {}) {
	        return new ExtractVideoThumbnailRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.thumbnail = source["thumbnail"];
	        this.mimeType = source["mimeType"];
	    }
	}
	export class RemoveVideoIntroReq {
	    filePath: string;
	    introDuration: number;
	
	    static createFrom(source: any = {}) {
	        return new RemoveVideoIntroReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.introDuration = source["introDuration"];
	    }
	}
	export class RemoveVideoIntroRes {
	    success: boolean;
	    message: string;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new RemoveVideoIntroRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.cost = source["cost"];
	    }
	}
	export class RotateVideoReq {
	    filePath: string;
	    angle: number;
	    clockwise: boolean;
	
	    static createFrom(source: any = {}) {
	        return new RotateVideoReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.angle = source["angle"];
	        this.clockwise = source["clockwise"];
	    }
	}
	export class RotateVideoRes {
	    success: boolean;
	    outputPath: string;
	    message: string;
	    cost: string;
	
	    static createFrom(source: any = {}) {
	        return new RotateVideoRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.outputPath = source["outputPath"];
	        this.message = source["message"];
	        this.cost = source["cost"];
	    }
	}
	export class ScanVideoDirReq {
	    dirPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanVideoDirReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirPath = source["dirPath"];
	    }
	}
	export class VideoExtInfo {
	    ext: string;
	    count: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoExtInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.ext = source["ext"];
	        this.count = source["count"];
	    }
	}
	export class ScanVideoDirRes {
	    success: boolean;
	    message: string;
	    totalCount: number;
	    extInfos: VideoExtInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ScanVideoDirRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.totalCount = source["totalCount"];
	        this.extInfos = this.convertValues(source["extInfos"], VideoExtInfo);
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
	export class ScanVideosReq {
	    dirPath: string;
	
	    static createFrom(source: any = {}) {
	        return new ScanVideosReq(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.dirPath = source["dirPath"];
	    }
	}
	export class VideoInfo {
	    filePath: string;
	    fileName: string;
	    fileSizeMB: number;
	    resolution: string;
	    width: number;
	    height: number;
	
	    static createFrom(source: any = {}) {
	        return new VideoInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.filePath = source["filePath"];
	        this.fileName = source["fileName"];
	        this.fileSizeMB = source["fileSizeMB"];
	        this.resolution = source["resolution"];
	        this.width = source["width"];
	        this.height = source["height"];
	    }
	}
	export class ScanVideosRes {
	    success: boolean;
	    message: string;
	    allVideos: VideoInfo[];
	    verticalVideos: VideoInfo[];
	
	    static createFrom(source: any = {}) {
	        return new ScanVideosRes(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.success = source["success"];
	        this.message = source["message"];
	        this.allVideos = this.convertValues(source["allVideos"], VideoInfo);
	        this.verticalVideos = this.convertValues(source["verticalVideos"], VideoInfo);
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

