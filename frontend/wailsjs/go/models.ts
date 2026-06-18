export namespace main {
	
	export class updateInfo {
	    has_update: boolean;
	    latest: string;
	    current: string;
	    download_url: string;
	
	    static createFrom(source: any = {}) {
	        return new updateInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.has_update = source["has_update"];
	        this.latest = source["latest"];
	        this.current = source["current"];
	        this.download_url = source["download_url"];
	    }
	}

}

export namespace model {
	
	export class CPUData {
	    model: string;
	    max_freq: string;
	    cores_threads: string;
	    cache_l3: string;
	    package_temp: number;
	    usage: number;
	    clock_speed: number;
	    vcore: number;
	    fan_speed: number;
	    power: number;
	    per_core_usage: number[];
	    mb_temp: number;
	    vrm_temp: number;
	
	    static createFrom(source: any = {}) {
	        return new CPUData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.model = source["model"];
	        this.max_freq = source["max_freq"];
	        this.cores_threads = source["cores_threads"];
	        this.cache_l3 = source["cache_l3"];
	        this.package_temp = source["package_temp"];
	        this.usage = source["usage"];
	        this.clock_speed = source["clock_speed"];
	        this.vcore = source["vcore"];
	        this.fan_speed = source["fan_speed"];
	        this.power = source["power"];
	        this.per_core_usage = source["per_core_usage"];
	        this.mb_temp = source["mb_temp"];
	        this.vrm_temp = source["vrm_temp"];
	    }
	}
	export class DiskIOData {
	    read_mbps: number;
	    write_mbps: number;
	
	    static createFrom(source: any = {}) {
	        return new DiskIOData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.read_mbps = source["read_mbps"];
	        this.write_mbps = source["write_mbps"];
	    }
	}
	export class NetworkData {
	    download_mbps: number;
	    upload_mbps: number;
	    connection_count: number;
	
	    static createFrom(source: any = {}) {
	        return new NetworkData(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.download_mbps = source["download_mbps"];
	        this.upload_mbps = source["upload_mbps"];
	        this.connection_count = source["connection_count"];
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
	    vram_type: string;
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
	        this.vram_type = source["vram_type"];
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
	    process_count: number;
	    thread_count: number;
	    ip_address: string;
	    display_info: string;
	
	    static createFrom(source: any = {}) {
	        return new SystemInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.time = source["time"];
	        this.date = source["date"];
	        this.uptime = source["uptime"];
	        this.process_count = source["process_count"];
	        this.thread_count = source["thread_count"];
	        this.ip_address = source["ip_address"];
	        this.display_info = source["display_info"];
	    }
	}
	export class DashboardData {
	    system: SystemInfo;
	    cpu: CPUData;
	    gpu: GPUData;
	    memory: MemoryData;
	    storage: StorageDrive[];
	    fans: FanInfo[];
	    total_power: number;
	    network: NetworkData;
	    disk_io: DiskIOData;
	
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
	        this.total_power = source["total_power"];
	        this.network = this.convertValues(source["network"], NetworkData);
	        this.disk_io = this.convertValues(source["disk_io"], DiskIOData);
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

