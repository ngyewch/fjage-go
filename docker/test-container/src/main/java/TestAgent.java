import org.arl.fjage.Agent;
import org.arl.fjage.param.ParameterMessageBehavior;

public class TestAgent extends Agent {
  @Override
  protected void init() {
    super.init();

    add(new ParameterMessageBehavior(TestParam.class));
  }

  public String[] getStrings() {
    return new String[] {"hello", "world"};
  }

  public int[] getInts() {
    return new int[] {12, -34, 56, -78, 89};
  }

  public float[] getFloats() {
    return new float[] {12.34f, -45.67f, 89f};
  }

  public double[] getDoubles() {
    return new double[] {12.34d, -45.67d, 89d};
  }
}
