<template>
  <div class="block gpu-block">
    <div class="block-header">
      <div class="block-title"><span class="dot gpu-dot"></span> GPU</div>
    </div>
    <div class="gpu-model">{{ gpu.model || '--' }}</div>
    <div class="gpu-specs">
      <span>{{ gpu.vram_spec || '--' }}</span>
    </div>
    <div class="metrics-grid">
      <MetricCard label="GPU TEMP" :value="gpu.temp" unit="°C" :max="100" bar-color="var(--purple)" :show-zero="false" />
      <MetricCard label="GPU USAGE" :value="gpu.usage" unit="%" :max="100" bar-color="var(--purple)" />
      <MetricCard label="GPU CLOCK" :value="gpu.clock" unit="MHz" :max="3000" bar-color="var(--purple)" :show-zero="false" />
      <MetricCard label="MEM USED" :value="gpu.mem_used" unit="MB" :max="gpu.mem_total || 16000" bar-color="var(--purple)" :show-zero="false" />
      <MetricCard label="FAN SPEED" :value="gpu.fan_speed" unit="RPM" :max="3000" bar-color="var(--purple)" :show-zero="false" />
      <MetricCard label="POWER" :value="gpu.power" unit="W" :max="300" bar-color="var(--purple)" :show-zero="false" />
    </div>
  </div>
</template>

<script>
import MetricCard from './MetricCard.vue'
export default {
  components: { MetricCard },
  props: { gpu: Object }
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

.gpu-block { border-top: 2px solid var(--purple); }

.block-header { margin-bottom: 8px; }
.block-title {
  font-size: 12px; font-weight: 700; color: var(--text-secondary);
  letter-spacing: 2px; display: flex; align-items: center; gap: 8px;
}
.dot { width: 8px; height: 8px; border-radius: 50%; display: inline-block; }
.gpu-dot { background: var(--purple); }

.gpu-model {
  font-size: 15px; font-weight: 600; color: var(--text-primary);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
  margin-bottom: 4px;
}

.gpu-specs {
  font-size: 12px; color: var(--text-secondary); margin-bottom: 12px;
}

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 10px;
  flex: 1;
}
</style>
