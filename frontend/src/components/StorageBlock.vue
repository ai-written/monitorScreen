<template>
  <div class="block stor-block">
    <div class="block-header">
      <div class="block-title"><span class="dot stor-dot"></span> STORAGE</div>
    </div>
    <div class="stor-list">
      <div v-if="!storage || storage.length === 0" class="stor-empty">No drives detected</div>
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
  </div>
</template>

<script>
export default {
  props: { storage: Array },
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
  padding: 16px 18px;
  display: flex;
  flex-direction: column;
}

.stor-block { border-top: 2px solid var(--amber); }

.block-header { margin-bottom: 8px; }
.block-title {
  font-size: 12px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.stor-dot { background: var(--amber); }

.stor-list { display: flex; flex-direction: column; gap: 12px; flex: 1; }
.stor-empty { color: var(--text-dim); font-size: 13px; text-align: center; padding: 20px 0; }

.stor-item { }
.stor-name {
  font-size: 13px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 2px;
}

.stor-brand {
  font-size: 11px; color: var(--text-dim);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 6px;
}

.stor-row { display: flex; justify-content: space-between; margin-bottom: 4px; }
.stor-temp { font-size: 12px; font-weight: 700; }
.stor-cap { font-size: 12px; color: var(--text-secondary); }

.progress-bar {
  height: 6px; background: var(--border); border-radius: 3px; overflow: hidden;
}
.progress-fill { height: 100%; border-radius: 3px; transition: width 0.5s; }
.stor-fill { background: linear-gradient(90deg, var(--amber), #fbbf24); }
</style>
