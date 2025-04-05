# Projections
# equirectangular
_A linear mapping between cartesian x/y and lat/lon._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| lat | 90.000000 | -90.000000 | ° | The latitude where 90 is North and -90 is South. |
| lon | -180.000000 | 180.000000 | ° | The longitude where -180 is West and 180 is East.. |

# mercator
_A cylindrical projection preserving angles, with significant area distortion near the poles._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| lat | 85.051129 | -85.051129 | ° | The latitude where 85.051° is the maximum and -85.051° is the minimum. |
| lon | -180.000000 | 180.000000 | ° | The longitude where -180 is West and 180 is East. |

# polar-to-rect
_Converts polar coordinates (angle and normalized radius) to rectangular (Cartesian) coordinates._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| θ | -3.141593 | 3.141593 | rad | Angular coordinate in radians (from -π to π). |
| r | 0.000000 | 1.000000 |  | Normalized radial coordinate (0 at center to 1 at maximum). |

# rect-to-polar
_Converts rectangular (Cartesian) coordinates to polar coordinates (angle and normalized radius)._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| θ | -3.141593 | 3.141593 | rad | Angular coordinate in radians (from -π to π). |
| r | 0.000000 | 1.000000 |  | Normalized radial coordinate (0 at center to 1 at maximum). |

# sinusoidal
_A pseudocylindrical projection that preserves area, with a non-linear x scaling based on the cosine of latitude._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| lat | 90.000000 | -90.000000 | ° | The latitude where 90 is North and -90 is South. |
| lon | -180.000000 | 180.000000 | ° | The longitude where -180 is West and 180 is East. |

# stereographic
_A projection that maps the sphere onto a plane using stereographic projection from the North Pole._

| Component | Min | Max | Unit | Desc |
| --- | --- | --- | --- | --- |
| lat | 90.000000 | -90.000000 | ° | The latitude where 90 is North and -90 is South. |
| lon | -180.000000 | 180.000000 | ° | The longitude where -180 is West and 180 is East. |


