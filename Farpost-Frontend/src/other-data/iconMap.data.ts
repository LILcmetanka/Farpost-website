import type { FC, SVGProps } from "react";
import HeatingIcon from "../assets/icons/heating.svg?react";
import HotWaterIcon from "../assets/icons/hot-water.svg?react";
import ColdWaterIcon from "../assets/icons/cold-water.svg?react";
import LightingIcon from "../assets/icons/lighting.svg?react";
import BuildingIcon from "../assets/icons/building.svg?react";
import CircleIcon from "../assets/icons/circle.svg?react";

export const ICONMAP: Record<string, FC<SVGProps<SVGSVGElement>>> = {
  heat: HeatingIcon,
  hot_water: HotWaterIcon,
  cold_water: ColdWaterIcon,
  electricity: LightingIcon,
  building: BuildingIcon,
  circle: CircleIcon,
};
