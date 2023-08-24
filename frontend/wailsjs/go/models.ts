export namespace src {
	
	export class Image {
	    id: string;
	    title: string;
	    image_url: string;
	    description: string;
	    category: string[];
	    tags: string[];
	
	    static createFrom(source: any = {}) {
	        return new Image(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.title = source["title"];
	        this.image_url = source["image_url"];
	        this.description = source["description"];
	        this.category = source["category"];
	        this.tags = source["tags"];
	    }
	}

}

