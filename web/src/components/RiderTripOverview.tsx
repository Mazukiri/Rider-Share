import { RouteFare, TripPreview, Driver } from "../types"
import { DriverList } from "./DriversList"
import { Card } from "./ui/card"
import { Button } from "./ui/button"
import { convertMetersToKilometers, convertSecondsToMinutes } from "../utils/math"
import { Skeleton } from "./ui/skeleton"
import { TripOverviewCard } from "./TripOverviewCard"
import { StripePaymentButton } from "./StripePaymentButton"
import { DriverCard } from "./DriverCard"
import { TripEvents, PaymentEventSessionCreatedData } from "../contracts"

interface TripOverviewProps {
  trip: TripPreview | null;
  status: TripEvents | null;
  assignedDriver?: Driver | null;
  paymentSession?: PaymentEventSessionCreatedData | null;
  onPackageSelect: (carPackage: RouteFare) => void;
  onCancel: () => void;
}