export namespace main {
	
	export class FileConversionInfo {
	    OriginalFileName: string;
	    OriginalFileSize: number;
	    NewFileName: string;
	    NewFileSize: number;
	    ConversionTime: number;
	    Thumbnail: number[];
	
	    static createFrom(source: any = {}) {
	        return new FileConversionInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.OriginalFileName = source["OriginalFileName"];
	        this.OriginalFileSize = source["OriginalFileSize"];
	        this.NewFileName = source["NewFileName"];
	        this.NewFileSize = source["NewFileSize"];
	        this.ConversionTime = source["ConversionTime"];
	        this.Thumbnail = source["Thumbnail"];
	    }
	}

}

