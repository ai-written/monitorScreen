<template>
  <div class="block stor-block">
    <div class="block-header">
      <div class="block-title"><span class="dot stor-dot"></span> 存储</div>
    </div>
    <div class="stor-list">
      <div v-if="!storage || storage.length === 0" class="stor-empty">未检测到硬盘</div>
      <div v-for="(d, i) in storage" :key="i" class="stor-item">
        <div class="stor-name">{{ d.name }}</div>
        <div v-if="d.brand" class="stor-brand">{{ d.brand }}</div>
        <div class="stor-row">
          <span class="stor-temp" :class="tempColor(d.temp)" v-if="d.temp > 0">{{ d.temp.toFixed(0) }}°C</span>
          <span class="stor-cap">{{ fmt(d.used_gb) }} / {{ fmt(d.total_gb) }} GB</span>
        </div>
        <div class="progress-bar">
          <div class="progress-fill stor-fill" :style="{ width: drvPct(d) + '%' }"></div>
        </div>
      </div>
    </div>
    <div class="disk-io" v-if="diskIO && (diskIO.read_mbps > 0 || diskIO.write_mbps > 0)">
      <div class="io-row">
        <span class="io-label">磁盘 IO</span>
        <span class="io-val io-read">R {{ fmt(diskIO.read_mbps) }} MB/s</span>
        <span class="io-val io-write">W {{ fmt(diskIO.write_mbps) }} MB/s</span>
      </div>
    </div>
  </div>
</template>

<script>
export default {
  props: { storage: Array, diskIO: Object },
  methods: {
    fmt(v) { return typeof v === 'number' ? v.toFixed(1) : '--' },
    drvPct(d) { return d.total_gb > 0 ? (d.used_gb / d.total_gb) * 100 : 0 },
    tempColor(t) {
      if (t > 60) return 'var(--red)'
      if (t > 45) return 'var(--amber)'
      return 'var(--green)'
    }
  }
}
</script>

<style scoped>
.block {
  background: var(--bg-card);
  border: 1px solid var(--border);
  border-radius: 10px;
  padding: 18px 20px;
  display: flex;
  flex-direction: column;
}

.stor-block { border-top: 2px solid var(--amber); max-height: 240px; display: flex; flex-direction: column; }

.block-header { margin-bottom: 8px; flex-shrink: 0; }
.block-title {
  font-size: 14px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.stor-dot { background: var(--amber); }

.stor-list {
  display: flex; flex-direction: column; gap: 10px; flex: 1;
  overflow-y: auto; min-height: 0;
  scrollbar-width: thin; scrollbar-color: var(--border) transparent;
}
.stor-list::-webkit-scrollbar { width: 4px; }
.stor-list::-webkit-scrollbar-track { background: transparent; }
.stor-list::-webkit-scrollbar-thumb { background: var(--border); border-radius: 2px; }
.stor-empty { color: var(--text-dim); font-size: 14px; text-align: center; padding: 20px 0; }

.stor-item { }
.stor-name {
  font-size: 14px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 2px;
}

.stor-brand {
  font-size: 12px; color: var(--text-dim);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 6px;
}

.stor-row { display: flex; justify-content: space-between; margin-bottom: 4px; }
.stor-temp { font-size: 13px; font-weight: 700; }
.stor-cap { font-size: 13px; color: var(--text-secondary); }

.progress-bar {
  height: 6px; background: var(--border); border-radius: 3px; overflow: hidden;
}
.progress-fill { height: 100%; border-radius: 3px; transition: width 0.5s; }
.stor-fill { background: linear-gradient(90deg, var(--amber), #fbbf24); }

.disk-io {
  margin-top: auto;
  padding-top: 8px;
  border-top: 1px solid var(--border);
  flex-shrink: 0;
}
.io-row {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 12px;
}
.io-label {
  color: var(--text-dim);
  letter-spacing: 1px;
}
.io-val { font-weight: 600; }
.io-read { color: var(--blue); }
.io-write { color: var(--amber); }
</style>
