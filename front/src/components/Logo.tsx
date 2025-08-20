import './Logo.css';

export default function Logo(): JSX.Element {
  return (
    <div className="dash-logo">
      {/* <img src={logo} alt="DashOPS - Beta" /> */}
      DashOPS <span>Beta</span>
    </div>
  );
}
