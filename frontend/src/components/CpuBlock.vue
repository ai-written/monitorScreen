<template>
  <div class="block cpu-block">
    <div class="block-header">
      <div class="block-title"><span class="dot cpu-dot"></span> CPU</div>
    </div>
    <div class="cpu-model">{{ cpu.model || '--' }}</div>
    <div class="cpu-specs">
      <span>{{ cpu.max_freq || '--' }}</span>
      <span class="sep">|</span>
      <span>{{ cpu.cores_threads || '--' }}</span>
      <template v-if="cpu.cache_l3">
        <span class="sep">|</span>
        <span>{{ cpu.cache_l3 }}</span>
      </template>
    </div>
    <div class="metrics-grid">
      <MetricCard label="CPU 温度" :value="cpu.package_temp" unit="°C" :max="100" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="CPU 使用率" :value="cpu.usage" unit="%" :max="100" bar-color="var(--cyan)" />
      <MetricCard label="主频" :value="cpu.clock_speed" unit="MHz" :max="6000" bar-color="var(--cyan)" />
      <MetricCard label="核心电压" :value="cpu.vcore" unit="V" :max="2" bar-color="var(--cyan)" :decimals="3" :show-zero="false" />
      <MetricCard label="主板温度" :value="cpu.mb_temp" unit="°C" :max="80" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="VRM 温度" :value="cpu.vrm_temp" unit="°C" :max="100" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="风扇转速" :value="cpu.fan_speed" unit="RPM" :max="2000" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="功耗" :value="cpu.power" unit="W" :max="200" bar-color="var(--cyan)" :show-zero="false" />
    </div>
    <div class="cores-section" v-if="cpu.per_core_usage && cpu.per_core_usage.length > 0">
      <div class="cores-label">各核心使用率</div>
      <div class="cores-bars">
        <div v-for="(pct, i) in cpu.per_core_usage" :key="i" class="core-item">
          <span class="core-num">{{ i }}</span>
          <div class="core-bar-track">
            <div class="core-bar-fill" :style="{ width: pct + '%', background: coreColor(pct) }"></div>
          </div>
          <span class="core-pct">{{ pct.toFixed(0) }}%</span>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import MetricCard from './MetricCard.vue'
export default {
  components: { MetricCard },
  props: { cpu: Object },
  methods: {
    coreColor(v) {
      if (v > 90) return 'var(--red)'
      if (v > 70) return 'var(--amber)'
      return 'var(--cyan)'
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
  overflow: hidden;
}

.cpu-block { border-top: 2px solid var(--cyan); }

.block-header { margin-bottom: 8px; }
.block-title {
  font-size: 14px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.cpu-dot { background: var(--cyan); }

.cpu-model {
  font-size: 17px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 4px;
}

.cpu-specs {
  font-size: 13px; color: var(--text-secondary); margin-bottom: 12px;
}
.sep { margin: 0 8px; color: var(--text-dim); }

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 6px;
  flex: 1;
}

.cores-section {
  margin-top: 12px;
  border-top: 1px solid var(--border);
  padding-top: 10px;
}
.cores-label { font-size: 11px; color: var(--text-dim); letter-spacing: 1px; margin-bottom: 6px; }
.cores-bars { display: flex; flex-direction: column; gap: 3px; max-height: 160px; overflow-y: auto; scrollbar-width: none; -ms-overflow-style: none; }
.cores-bars::-webkit-scrollbar { display: none; }
.core-item { display: flex; align-items: center; gap: 6px; }
.core-num { font-size: 11px; color: var(--text-dim); width: 16px; text-align: right; flex-shrink: 0; }
.core-bar-track { flex: 1; height: 6px; background: var(--border); border-radius: 3px; overflow: hidden; }
.core-bar-fill { height: 100%; border-radius: 3px; transition: width 0.4s; min-width: 2px; }
.core-pct { font-size: 11px; color: var(--text-secondary); width: 28px; flex-shrink: 0; }
</style>
