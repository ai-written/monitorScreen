export namespace model {
	
	export class CPUData {
	    model: string;
	    max_freq: string;
	    cores_threads: string;
	    package_temp: number;
	    usage: number;
	    clock_speed: number;
	    vcore: number;
	    fan_speed: number;
	    power: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.max_freq = source["max_freq"];
	        this.cores_threads = source["cores_threads"];
	        this.package_temp = source["package_temp"];
	        this.usage = source["usage"];
	        this.clock_speed = source["clock_speed"];
	        this.vcore = source["vcore"];
	        this.fan_speed = source["fan_speed"];
	        this.power = source["power"];
	    }
	}
	export class FanInfo {
	    name: string;
	    rpm: number;
	    percent: number;
	
	    static createFrom(source: any = {}) {
	        return new FanInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.rpm = source["rpm"];
	        this.percent = source["percent"];
	    }
	}
	export class StorageDrive {
	    name: string;
	    brand: string;
	    temp: number;
	    used_gb: number;
	    total_gb: number;
	
	    static createFrom(source: any = {}) {
	        return new StorageDrive(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.brand = source["brand"];
	        this.temp = source["temp"];
	        this.used_gb = source["used_gb"];
	        this.total_gb = source["total_gb"];
	    }
	}
	export class MemoryData {
	    used_gb: number;
	    total_gb: number;
	    type: string;
	    frequency: string;
	    channel: string;
	    brand: string;
	
	    static createFrom(source: any = {}) {
	        return new MemoryData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.used_gb = source["used_gb"];
	        this.total_gb = source["total_gb"];
	        this.type = source["type"];
	        this.frequency = source["frequency"];
	        this.channel = source["channel"];
	        this.brand = source["brand"];
	    }
	}
	export class GPUData {
	    model: string;
	    vram_spec: string;
	    pcie_version: string;
	    temp: number;
	    usage: number;
	    clock: number;
	    mem_used: number;
	    mem_total: number;
	    fan_speed: number;
	    power: number;
	
	    static createFrom(source: any = {}) {
	        return new GPUData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.vram_spec = source["vram_spec"];
	        this.pcie_version = source["pcie_version"];
	        this.temp = source["temp"];
	        this.usage = source["usage"];
	        this.clock = source["clock"];
	        this.mem_used = source["mem_used"];
	        this.mem_total = source["mem_total"];
	        this.fan_speed = source["fan_speed"];
	        this.power = source["power"];
	    }
	}
	export class SystemInfo {
	    time: string;
	    date: string;
	    uptime: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.date = source["date"];
	        this.uptime = source["uptime"];
	    }
	}
	export class DashboardData {
	    system: SystemInfo;
	    cpu: CPUData;
	    gpu: GPUData;
	    memory: MemoryData;
	    storage: StorageDrive[];
	    fans: FanInfo[];
	
	    static createFrom(source: any = {}) {
	        return new DashboardData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.system = this.convertValues(source["system"], SystemInfo);
	        this.cpu = this.convertValues(source["cpu"], CPUData);
	        this.gpu = this.convertValues(source["gpu"], GPUData);
	        this.memory = this.convertValues(source["memory"], MemoryData);
	        this.storage = this.convertValues(source["storage"], StorageDrive);
	        this.fans = this.convertValues(source["fans"], FanInfo);
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

