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
    </div>
    <div class="metrics-grid">
      <MetricCard label="PACKAGE TEMP" :value="cpu.package_temp" unit="°C" :max="100" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="CPU USAGE" :value="cpu.usage" unit="%" :max="100" bar-color="var(--cyan)" />
      <MetricCard label="CLOCK SPEED" :value="cpu.clock_speed" unit="MHz" :max="6000" bar-color="var(--cyan)" />
      <MetricCard label="VCORE" :value="cpu.vcore" unit="V" :max="2" bar-color="var(--cyan)" :decimals="3" :show-zero="false" />
      <MetricCard label="FAN SPEED" :value="cpu.fan_speed" unit="RPM" :max="2000" bar-color="var(--cyan)" :show-zero="false" />
      <MetricCard label="POWER" :value="cpu.power" unit="W" :max="200" bar-color="var(--cyan)" :show-zero="false" />
    </div>
  </div>
</template>

<script>
import MetricCard from './MetricCard.vue'
export default {
  components: { MetricCard },
  props: { cpu: Object }
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
  overflow: hidden;
}

.cpu-block { border-top: 2px solid var(--cyan); }

.block-header { margin-bottom: 8px; }
.block-title {
  font-size: 12px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.cpu-dot { background: var(--cyan); }

.cpu-model {
  font-size: 15px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 4px;
}

.cpu-specs {
  font-size: 12px; color: var(--text-secondary); margin-bottom: 12px;
}
.sep { margin: 0 8px; color: var(--text-dim); }

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  flex: 1;
}
</style>
