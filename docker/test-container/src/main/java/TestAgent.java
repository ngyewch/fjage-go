import org.arl.fjage.Agent;
import org.arl.fjage.param.ParameterMessageBehavior;

public class TestAgent extends Agent {
  private String[] strings = new String[] {"hello, world"};
  private int[] ints = new int[] {12, -34, 56, -78, 89};
  private float[] floats = new float[] {12.34f, -45.67f, 89f};
  private double[] doubles = new double[] {12.34d, -45.67d, 89d};

  @Override
  protected void init() {
    super.init();

    add(new ParameterMessageBehavior(TestParam.class));
  }

  public String[] getStrings() {
    return strings;
  }

  public void setStrings(String[] values) {
    strings = values;
  }

  public int[] getInts() {
    return ints;
  }

  public void setInts(int[] values) {
    this.ints = values;
  }

  public float[] getFloats() {
    return floats;
  }

  public void setFloats(float[] values) {
    this.floats = values;
  }

  public double[] getDoubles() {
    return doubles;
  }

  public void setDoubles(double[] values) {
    this.doubles = values;
  }
}
